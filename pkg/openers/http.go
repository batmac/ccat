//go:build !fileonly
// +build !fileonly

package openers

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"

	"braces.dev/errtrace"
	"github.com/batmac/ccat/pkg/globalctx"
	"github.com/batmac/ccat/pkg/log"
)

var (
	httpOpenerName        = "http"
	httpOpenerDescription = "get URL via HTTP(S)"
)

type httpOpener struct {
	name, description string
}

func init() {
	register(&httpOpener{
		name:        httpOpenerName,
		description: httpOpenerDescription,
	})
	// http.DefaultClient.Timeout = 10 * time.Second
}

func (f httpOpener) Name() string {
	return f.name
}

func (f httpOpener) Description() string {
	return f.description
}

func (f httpOpener) Open(s string, _ bool) (io.ReadCloser, error) {
	flag := globalctx.GetBool("insecure")
	tr := http.DefaultTransport.(*http.Transport)
	// log.Debugf("flag=%v, tr=%#v\n", flag, tr)
	//#nosec
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: flag}

	//#nosec
	resp, err := http.Get(s)
	if err != nil {
		log.Println(err)
		return nil, errtrace.Wrap(err)
	}
	nrc := resp.Body

	return nrc, nil
}

func (f httpOpener) Evaluate(s string) float32 {
	// log.Debugf("Evaluating %s...\n", s)
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return 0.9
	}
	return 0
}
