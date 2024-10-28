//go:build !fileonly
// +build !fileonly

package openers

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/psanford/wormhole-william/wormhole"

	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/utils"
)

var (
	wormholeOpenerName        = "wormhole"
	wormholeOpenerDescription = "get text, file or zipped dir via a wormhole code (wh://<code> or wormhole://<code>)"
)

type wormholeOpener struct {
	name, description string
}

func init() {
	register(&wormholeOpener{
		name:        wormholeOpenerName,
		description: wormholeOpenerDescription,
	})
}

func (f wormholeOpener) Name() string {
	return f.name
}

func (f wormholeOpener) Description() string {
	return f.description
}

func (f wormholeOpener) Open(code string, _ bool) (io.ReadCloser, error) {
	if strings.HasPrefix(code, "wormhole://") {
		code = strings.TrimPrefix(code, "wormhole://")
	} else if strings.HasPrefix(code, "wh://") {
		code = strings.TrimPrefix(code, "wh://")
	}

	var c wormhole.Client

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	fileInfo, err := c.Receive(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("got msg: %+v\n", fileInfo)

	return utils.NewReadCloser(fileInfo, func() error {
		cancel()
		return nil
	}), nil
}

func (f wormholeOpener) Evaluate(s string) float32 {
	if strings.HasPrefix(s, "wormhole://") || strings.HasPrefix(s, "wh://") {
		return 0.9
	}
	return 0
}
