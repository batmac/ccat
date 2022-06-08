//go:build !windows
// +build !windows

package openers

import (
	"ccat/log"
	"io"

	curl "github.com/andelf/go-curl"
)

var easy *curl.CURL = nil

var curlOpenerName = "curl"
var curlOpenerDescription = "get URL via go-curl (C bindings)"

type curlOpener struct {
	name, description string
}

func init() {
	register(&curlOpener{
		name:        curlOpenerName,
		description: curlOpenerDescription,
	})
}

func easyHandlerInit() {
	easy = curl.EasyInit()
	easy.Setopt(curl.OPT_VERBOSE, false)
	easy.Setopt(curl.OPT_WRITEFUNCTION, func(ptr []byte, userdata interface{}) bool {
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
func (f curlOpener) Open(s string, _ bool) (io.ReadCloser, error) {

	r, w := io.Pipe()
	go func() {
		// we don't cleanup curl stuff when ending because we don't care (we only use one)
		log.Debugln(" curl goroutine started")
		//curl.GlobalInit(curl.GLOBAL_DEFAULT)
		//defer curl.GlobalCleanup()
		if easy == nil {
			easyHandlerInit()
		}
		//defer easy.Cleanup()

		easy.Setopt(curl.OPT_URL, s)

		easy.Setopt(curl.OPT_WRITEDATA, w)

		if err := easy.Perform(); err != nil {
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
	//log.Debugf("Evaluating %s...\n", s)
	return 0.1
}
