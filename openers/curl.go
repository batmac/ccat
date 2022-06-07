package openers

import (
	"ccat/log"
	"io"

	curl "github.com/andelf/go-curl"
)

var curlOpenerName = "curl"
var curlOpenerDescription = "get URL via go-curl"

type curlOpener struct {
	name, description string
}

func init() {
	register(&curlOpener{
		name:        curlOpenerName,
		description: curlOpenerDescription,
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
		log.Debugln(" curl goroutine started")
		curl.GlobalInit(curl.GLOBAL_DEFAULT)
		defer curl.GlobalCleanup()
		easy := curl.EasyInit()
		defer easy.Cleanup()

		easy.Setopt(curl.OPT_URL, s)

		easy.Setopt(curl.OPT_WRITEFUNCTION, func(ptr []byte, userdata interface{}) bool {
			pipe := userdata.(*io.PipeWriter)
			if _, err := pipe.Write(ptr); err != nil {
				return false
			}
			return true
		})

		easy.Setopt(curl.OPT_WRITEDATA, w)
		easy.Setopt(curl.OPT_VERBOSE, false)
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
	//log.Debugf("Evaluating %s...\n", s)
	return 0.1
}
