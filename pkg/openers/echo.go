//go:build !fileonly
// +build !fileonly

package openers

import (
	"fmt"
	"io"
	"strings"
)

var (
	echoOpenerName        = "echo"
	echoOpenerDescription = "echo the string given"
)

type echoOpener struct {
	name, description string
}

func init() {
	register(&echoOpener{
		name:        echoOpenerName,
		description: echoOpenerDescription,
	})
}

func (f echoOpener) Name() string {
	return f.name
}

func (f echoOpener) Description() string {
	return f.description
}

func (f echoOpener) Open(s string, _ bool) (io.ReadCloser, error) {
	var found bool
	_, s, found = strings.Cut(s, "://")
	if !found {
		return nil, fmt.Errorf("no protocol given")
	}
	// remove quotes from start and end
	s = strings.Trim(s, "\"'")

	if len(s) == 0 {
		return nil, fmt.Errorf("no data given")
	}

	datareader := io.NopCloser(strings.NewReader(s))
	return datareader, nil
}

func (f echoOpener) Evaluate(s string) float32 {
	// log.Debugf("Evaluating %s...\n", s)
	if strings.HasPrefix(s, "echo://") {
		return 0.9
	}
	return 0
}
