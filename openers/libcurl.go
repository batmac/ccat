//go:build cgo && libcurl
// +build cgo,libcurl

package openers

import (
	"ccat/log"
	"ccat/utils"
	"io"
	"strings"

	curl "github.com/andelf/go-curl"
)

var curlOpenerName = "curl"
var curlOpenerDescription = "get URL via libcurl bindings\n           " +
	curl.Version() + "\n           protocols: " +
	strings.Join(curl.VersionInfo(0).Protocols, ",")

type curlOpener struct {
	easy              *curl.CURL
	name, description string
}

func init() {
	register(&curlOpener{
		name:        curlOpenerName,
		description: curlOpenerDescription,
	})
}

func (f *curlOpener) easyHandlerInit() {
	// we don't cleanup curl stuff when ending because we don't care (we only use one)

	//curl.GlobalInit(curl.GLOBAL_DEFAULT)
	//defer curl.GlobalCleanup()
	f.easy = curl.EasyInit()
	f.easy.Setopt(curl.OPT_VERBOSE, false)
	f.easy.Setopt(curl.OPT_TIMEOUT, 10)
	f.easy.Setopt(curl.OPT_WRITEFUNCTION, func(ptr []byte, userdata interface{}) bool {
		pipe := userdata.(*io.PipeWriter)
		if _, err := pipe.Write(ptr); err != nil {
			return false
		}
		return true
	})
}

func (f curlOpener) Name() string {
	return f.name
}
func (f curlOpener) Description() string {
	return f.description
}
func (f *curlOpener) Open(s string, _ bool) (io.ReadCloser, error) {

	r, w := io.Pipe()
	go func() {
		log.Debugln(" curl goroutine started")

		if f.easy == nil {
			f.easyHandlerInit()
		}
		//defer easy.Cleanup()

		f.easy.Setopt(curl.OPT_URL, s)

		f.easy.Setopt(curl.OPT_WRITEDATA, w)

		if err := f.easy.Perform(); err != nil {
			println(" curl ERROR", err.Error())
			w.CloseWithError(err)
		}
		w.Close()
		log.Debugln(" curl goroutine ended")
	}()

	return r, nil
}

func (f curlOpener) Evaluate(s string) float32 {
	// https://everything.curl.dev/protocols/curl
	// The latest curl (as of this writing) supports these protocols:
	// DICT, FILE, FTP, FTPS, GOPHER, GOPHERS, HTTP, HTTPS, IMAP, IMAPS, LDAP, LDAPS,
	// MQTT, POP3, POP3S, RTMP, RTSP, SCP, SFTP, SMB, SMBS, SMTP, SMTPS, TELNET, TFTP
	before, _, found := strings.Cut(s, "://")
	//log.Printf("before=%s found=%v s=%v", before, found, s)
	if found && utils.StringInSlice(before, curl.VersionInfo(0).Protocols) {
		return 0.1
	}
	return 0
}
