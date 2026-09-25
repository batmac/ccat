//go:build cgo && libcurl && !fileonly

package openers

import (
	"io"
	"strings"
	"time"

	"github.com/batmac/ccat/pkg/globalctx"
	"github.com/batmac/ccat/pkg/log"

	curl "github.com/batmac/go-curl"
)

var (
	curlOpenerDescription = "get URL via libcurl bindings\n           " +
		curl.Version() + "\n           protocols: " +
		strings.Join(curl.VersionInfo(0).Protocols, ",")
)

type curlOpener struct {
	easy              *curl.CURL
	name, description string
}

func init() {
	register(&curlOpener{
		easy:        nil,
		name:        curlOpenerName,
		description: curlOpenerDescription,
	})
}

func (f *curlOpener) easyHandlerInit() {
	// we don't cleanup curl stuff when ending because we don't care (we only use one)

	// curl.GlobalInit(curl.GLOBAL_DEFAULT)
	// defer curl.GlobalCleanup()
	f.easy = curl.EasyInit()
	if err := f.easy.Setopt(curl.OPT_VERBOSE, false); err != nil {
		log.Println(" curl ERROR", err.Error())
	}
	if err := f.easy.Setopt(curl.OPT_FOLLOWLOCATION, true); err != nil {
		log.Println(" curl ERROR", err.Error())
	}
	if err := f.easy.Setopt(curl.OPT_MAXREDIRS, 10); err != nil {
		log.Println(" curl ERROR", err.Error())
	}
	if err := f.easy.Setopt(curl.OPT_CONNECTTIMEOUT, 10); err != nil {
		log.Println(" curl ERROR", err.Error())
	}
	if err := f.easy.Setopt(curl.OPT_WRITEFUNCTION, func(ptr []byte, userdata any) bool {
		pipe := userdata.(*io.PipeWriter)
		if _, err := pipe.Write(ptr); err != nil {
			return false
		}
		return true
	}); err != nil {
		log.Println(" curl ERROR", err.Error())
	}

	step := time.Now().Unix()
	dlstep := 0.0
	if err := f.easy.Setopt(curl.OPT_NOPROGRESS, false); err != nil {
		log.Println(" curl ERROR", err.Error())
	}
	if err := f.easy.Setopt(curl.OPT_PROGRESSFUNCTION, func(dltotal, dlnow, ultotal, ulnow float64, _ any) bool {
		if time.Now().Unix()-step > 2 {
			log.Debugf("downloaded: %3.2f%%, speed: %.1fKiB/s \r", dlnow/dltotal*100, (dlnow-dlstep)/1000/float64((time.Now().Unix()-step)))
			step = time.Now().Unix()
			dlstep = dlnow
		}
		return true
	}); err != nil {
		log.Println(" curl ERROR", err.Error())
	}
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
		// defer easy.Cleanup()

		s = tryTransformURL(s)

		flag := globalctx.GetBool("insecure")
		if flag {
			log.Debugln(" curl insecure enabled!")
			if err := f.easy.Setopt(curl.OPT_SSL_VERIFYHOST, 0); err != nil {
				log.Println(" curl ERROR", err.Error())
			}
			if err := f.easy.Setopt(curl.OPT_SSL_VERIFYPEER, 0); err != nil {
				log.Println(" curl ERROR", err.Error())
			}
		} else {
			log.Debugln(" curl SECURE only.")
			if err := f.easy.Setopt(curl.OPT_SSL_VERIFYHOST, 2); err != nil {
				log.Println(" curl ERROR", err.Error())
			}
			if err := f.easy.Setopt(curl.OPT_SSL_VERIFYPEER, 1); err != nil {
				log.Println(" curl ERROR", err.Error())
			}
		}

		if err := f.easy.Setopt(curl.OPT_URL, s); err != nil {
			log.Println(" curl ERROR", err.Error())
		}
		if err := f.easy.Setopt(curl.OPT_WRITEDATA, w); err != nil {
			log.Println(" curl ERROR", err.Error())
		}

		if err := f.easy.Perform(); err != nil {
			log.Println(" curl ERROR", err.Error())
			if err := w.CloseWithError(err); err != nil {
				log.Println(" curl ERROR", err.Error())
			}
		}
		if err := w.Close(); err != nil {
			log.Println(" curl ERROR", err.Error())
		}
		log.Debugln(" curl goroutine ended")
	}()

	return r, nil
}

func (f curlOpener) Evaluate(s string) float32 {
	return curlEvaluate(s, curl.VersionInfo(0).Protocols)
}
