//go:build !fileonly
// +build !fileonly

package openers

import (
	"context"
	"io"
	"strings"

	"git.sr.ht/~adnano/go-gemini"
	"git.sr.ht/~adnano/go-gemini/tofu"

	"braces.dev/errtrace"
	"github.com/batmac/ccat/pkg/log"
)

var (
	geminiOpenerName        = "gemini"
	geminiOpenerDescription = "get URL via Gemini"

	knownHosts = tofu.KnownHosts{}
)

type geminiOpener struct {
	name, description string
}

func init() {
	register(&geminiOpener{
		name:        geminiOpenerName,
		description: geminiOpenerDescription,
	})
}

func (f geminiOpener) Name() string {
	return f.name
}

func (f geminiOpener) Description() string {
	return f.description
}

func (f geminiOpener) Open(s string, _ bool) (io.ReadCloser, error) {
	log.Debugf("gemini knownHosts: %v\n", knownHosts.Entries())
	client := &gemini.Client{
		TrustCertificate: knownHosts.TOFU,
	}
	ctx := context.Background()
	resp, err := client.Get(ctx, s)
	if err != nil {
		// handle error
		log.Println(err)
		return nil, errtrace.Wrap(err)
	}

	return resp.Body, nil
}

func (f geminiOpener) Evaluate(s string) float32 {
	// log.Debugf("Evaluating %s...\n", s)
	if strings.HasPrefix(s, "gemini://") {
		return 0.9
	}
	return 0
}
