//go:build darwin && libcurl_purego && !fileonly && !(cgo && libcurl)

package openers

// libcurl opener without cgo: libcurl is loaded at runtime with purego.
// see curlSetopt for the call to the variadic curl_easy_setopt.

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/batmac/ccat/pkg/globalctx"
	"github.com/batmac/ccat/pkg/log"
	"github.com/ebitengine/purego"
)

// libcurl ships with macOS (in the dyld shared cache), nothing has to be installed.
const libcurlPath = "/usr/lib/libcurl.4.dylib"

// values from curl/curl.h
const (
	curlGlobalDefault = 3

	curloptNoProgress       = 43
	curloptFollowLocation   = 52
	curloptSSLVerifyPeer    = 64
	curloptMaxRedirs        = 68
	curloptConnectTimeout   = 78
	curloptSSLVerifyHost    = 81
	curloptNoSignal         = 99
	curloptWriteData        = 10001
	curloptURL              = 10002
	curloptXferInfoData     = 10057
	curloptWriteFunction    = 20011
	curloptXferInfoFunction = 20219
)

// curlVersionInfo mirrors the CURLVERSION_FIRST part of struct curl_version_info_data (LP64).
type curlVersionInfo struct {
	age           int32
	version       *byte
	versionNum    uint32
	host          *byte
	features      int32
	sslVersion    *byte
	sslVersionNum int64
	libzVersion   *byte
	protocols     **byte
}

type libcurl struct {
	version   string
	protocols []string

	easyInit     func() uintptr
	easyCleanup  func(uintptr)
	easyPerform  func(uintptr) int32
	easyStrerror func(int32) string

	// CURLcode curl_easy_setopt(CURL *handle, CURLoption option, ...);
	setopt uintptr

	writeCb, xferInfoCb uintptr
}

// a transfer is referenced by the callbacks through its id (passed as userdata),
// Go pointers can't be kept by C code.
type curlTransfer struct {
	w       *io.PipeWriter
	step    time.Time
	stepNow int64
}

var (
	curlTransfers  sync.Map // uintptr -> *curlTransfer
	curlTransferID atomic.Uintptr

	loadLibcurl = sync.OnceValues(newLibcurl)
)

type curlPuregoOpener struct {
	name string
}

func init() {
	register(&curlPuregoOpener{
		name: curlOpenerName,
	})
}

func (f curlPuregoOpener) Name() string {
	return f.name
}

func (f curlPuregoOpener) Description() string {
	c, err := loadLibcurl()
	if err != nil {
		return "get URL via libcurl (purego), unavailable: " + err.Error()
	}
	return "get URL via libcurl (purego)\n           " +
		c.version + "\n           protocols: " +
		strings.Join(c.protocols, ",")
}

func (f curlPuregoOpener) Evaluate(s string) float32 {
	// don't load libcurl for local files
	if !strings.Contains(s, "://") {
		return 0
	}
	c, err := loadLibcurl()
	if err != nil {
		log.Debugln(" curl unavailable:", err)
		return 0
	}
	return curlEvaluate(s, c.protocols)
}

func (f *curlPuregoOpener) Open(s string, _ bool) (io.ReadCloser, error) {
	c, err := loadLibcurl()
	if err != nil {
		return nil, err
	}

	h := c.easyInit()
	if h == 0 {
		return nil, errors.New("curl_easy_init failed")
	}

	r, w := io.Pipe()
	id := curlTransferID.Add(1)
	if err := c.setup(h, tryTransformURL(s), id); err != nil {
		c.easyCleanup(h)
		return nil, err
	}
	curlTransfers.Store(id, &curlTransfer{w: w, step: time.Now()})

	go func() {
		log.Debugln(" curl goroutine started")
		defer func() {
			curlTransfers.Delete(id)
			c.easyCleanup(h)
			log.Debugln(" curl goroutine ended")
		}()

		if rc := c.easyPerform(h); rc != 0 {
			err := fmt.Errorf("curl: %s", c.easyStrerror(rc))
			log.Println(" curl ERROR", err.Error())
			_ = w.CloseWithError(err)
			return
		}
		_ = w.Close()
	}()

	return r, nil
}

func (c *libcurl) setup(h uintptr, url string, id uintptr) error {
	verifyPeer, verifyHost := int64(1), int64(2)
	if globalctx.GetBool("insecure") {
		log.Debugln(" curl insecure enabled!")
		verifyPeer, verifyHost = 0, 0
	} else {
		log.Debugln(" curl SECURE only.")
	}

	for _, o := range []struct {
		opt int32
		val int64
	}{
		// libcurl must not use signals in a multithreaded (Go) program
		{curloptNoSignal, 1},
		{curloptFollowLocation, 1},
		{curloptMaxRedirs, 10},
		{curloptConnectTimeout, 10},
		{curloptNoProgress, 0},
		{curloptSSLVerifyPeer, verifyPeer},
		{curloptSSLVerifyHost, verifyHost},
	} {
		if err := c.setoptValue(h, o.opt, uintptr(o.val)); err != nil {
			return err
		}
	}

	for _, o := range []struct {
		opt int32
		val uintptr
	}{
		{curloptWriteFunction, c.writeCb},
		{curloptWriteData, id},
		{curloptXferInfoFunction, c.xferInfoCb},
		{curloptXferInfoData, id},
	} {
		if err := c.setoptValue(h, o.opt, o.val); err != nil {
			return err
		}
	}

	// libcurl copies string options
	cURL := append([]byte(url), 0)
	err := c.setoptValue(h, curloptURL, uintptr(unsafe.Pointer(&cURL[0])))
	runtime.KeepAlive(cURL)
	return err
}

func (c *libcurl) setoptValue(h uintptr, opt int32, val uintptr) error {
	if rc := curlSetopt(c.setopt, h, opt, val); rc != 0 {
		return fmt.Errorf("curl_easy_setopt(%d): %s", opt, c.easyStrerror(rc))
	}
	return nil
}

// curlSetopt calls the variadic curl_easy_setopt(handle, option, val) with a
// long or a pointer as val. purego doesn't support C variadic functions, so val
// is passed twice, to be where the callee reads its variadic argument:
//   - as the 3rd argument, in a register: SysV amd64 (purego sets AL=0, no
//     vector registers used) and standard AAPCS64,
//   - as the 9th argument: purego.SyscallN on arm64 puts the 9th argument in the
//     first stack slot (0(RSP)), which is where the Apple arm64 ABI passes the
//     variadic arguments.
//
// The extra arguments are ignored by the callee (the caller owns the stack).
func curlSetopt(setopt, h uintptr, opt int32, val uintptr) int32 {
	rc, _, _ := purego.SyscallN(setopt, h, uintptr(opt), val, 0, 0, 0, 0, 0, val)
	return int32(rc)
}

// size_t write_callback(char *ptr, size_t size, size_t nmemb, void *userdata);
func curlWrite(ptr *byte, size, nmemb, userdata uintptr) uintptr {
	n := size * nmemb
	t, ok := curlTransfers.Load(userdata)
	if !ok || n == 0 {
		return 0
	}
	// io.Pipe copies before returning, ptr isn't referenced afterwards.
	// returning less than n makes libcurl abort with CURLE_WRITE_ERROR
	if _, err := t.(*curlTransfer).w.Write(unsafe.Slice(ptr, n)); err != nil {
		return 0
	}
	return n
}

// int xferinfo_callback(void *clientp, curl_off_t dltotal, curl_off_t dlnow, curl_off_t ultotal, curl_off_t ulnow);
func curlXferInfo(clientp uintptr, dltotal, dlnow, _, _ int64) uintptr {
	v, ok := curlTransfers.Load(clientp)
	if !ok {
		return 0
	}
	t := v.(*curlTransfer)
	if elapsed := time.Since(t.step); elapsed > 2*time.Second {
		if dltotal > 0 {
			log.Debugf("downloaded: %3.2f%%, speed: %.1fKiB/s \r", float64(dlnow)/float64(dltotal)*100, float64(dlnow-t.stepNow)/1000/elapsed.Seconds())
		}
		t.step = time.Now()
		t.stepNow = dlnow
	}
	return 0
}

func newLibcurl() (*libcurl, error) {
	lib, err := purego.Dlopen(libcurlPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return nil, err
	}

	c := &libcurl{}
	var (
		globalInit  func(int64) int32
		version     func() string
		versionInfo func(int32) *curlVersionInfo
	)
	for name, fptr := range map[string]any{
		"curl_global_init":   &globalInit,
		"curl_version":       &version,
		"curl_version_info":  &versionInfo,
		"curl_easy_init":     &c.easyInit,
		"curl_easy_cleanup":  &c.easyCleanup,
		"curl_easy_perform":  &c.easyPerform,
		"curl_easy_strerror": &c.easyStrerror,
	} {
		addr, err := purego.Dlsym(lib, name)
		if err != nil {
			return nil, err
		}
		purego.RegisterFunc(fptr, addr)
	}

	if c.setopt, err = purego.Dlsym(lib, "curl_easy_setopt"); err != nil {
		return nil, err
	}

	if rc := globalInit(curlGlobalDefault); rc != 0 {
		return nil, fmt.Errorf("curl_global_init failed (%d)", rc)
	}

	c.version = version()
	if info := versionInfo(0); info != nil {
		for p := info.protocols; p != nil && *p != nil; p = (**byte)(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(p))) {
			c.protocols = append(c.protocols, cString(*p))
		}
	}

	// created once: purego callbacks are never released
	c.writeCb = purego.NewCallback(curlWrite)
	c.xferInfoCb = purego.NewCallback(curlXferInfo)

	return c, nil
}

func cString(p *byte) string {
	n := 0
	for *(*byte)(unsafe.Add(unsafe.Pointer(p), n)) != 0 {
		n++
	}
	return string(unsafe.Slice(p, n))
}
