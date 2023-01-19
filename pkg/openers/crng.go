//go:build !fileonly
// +build !fileonly

package openers

import (
	"crypto/rand"
	"io"
	"strings"

	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/utils"
	"github.com/docker/go-units"
)

var (
	crngOpenerName        = "crng"
	crngOpenerDescription = "get data from crypto/rand (accept a size limit as parameter)"
)

type crngOpener struct {
	name, description string
}

func init() {
	register(&crngOpener{
		name:        crngOpenerName,
		description: crngOpenerDescription,
	})
}

func (f crngOpener) Name() string {
	return f.name
}

func (f crngOpener) Description() string {
	return f.description
}

func (f crngOpener) Open(s string, _ bool) (io.ReadCloser, error) {
	s = strings.TrimPrefix(s, "crng://")

	R := rand.Reader
	var limit int64
	var err error

	if len(s) != 0 {
		limit, err = units.FromHumanSize(s)
		if err != nil {
			return nil, err
		}
		log.Debugln("limiting to ", limit, " bytes")
		R = io.LimitReader(rand.Reader, limit)
	}

	return utils.NewReadCloser(R, func() error { return nil }), nil
}

func (f crngOpener) Evaluate(s string) float32 {
	// log.Debugf("Evaluating %s...\n", s)
	if strings.HasPrefix(s, "crng://") {
		return 0.9
	}
	return 0
}
