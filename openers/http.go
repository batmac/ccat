package openers

import (
	"ccat/log"
	"io"
	"net/http"
	"strings"
)

var httpOpenerName = "http"
var httpOpenerDescription = "get URL via HTTP(S)"

type httpOpener struct {
	name, description string
}

func init() {
	register(&httpOpener{
		name:        httpOpenerName,
		description: httpOpenerDescription,
	})
}

func (f httpOpener) Name() string {
	return f.name
}
func (f httpOpener) Description() string {
	return f.description
}
func (f httpOpener) Open(s string, _ bool) (io.ReadCloser, error) {
	resp, err := http.Get(s)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	nrc := NewReadCloser(resp.Body, func() error {

		return nil
	})

	return nrc, nil
}

func (f httpOpener) Evaluate(s string) float32 {
	//log.Debugf("Evaluating %s...\n", s)
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return 0.9
	}
	return 0
}
