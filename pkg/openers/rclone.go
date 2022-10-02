//go:build !fileonly
// +build !fileonly

package openers

import (
	"context"
	"io"
	"strings"

	"github.com/batmac/ccat/pkg/log"

	// load all backends
	_ "github.com/rclone/rclone/backend/all"
	"github.com/rclone/rclone/cmd"
	"github.com/rclone/rclone/fs/config/configfile"
	"github.com/rclone/rclone/fs/operations"
)

var (
	rcloneOpenerName        = "rclone"
	rcloneOpenerDescription = "get data from rclone config"
)

type rcloneOpener struct {
	name, description string
}

func init() {
	register(&rcloneOpener{
		name:        rcloneOpenerName,
		description: rcloneOpenerDescription,
	})
}

func (f rcloneOpener) Name() string {
	return f.name
}

func (f rcloneOpener) Description() string {
	return f.description
}

func (f rcloneOpener) Open(s string, _ bool) (io.ReadCloser, error) {
	path := strings.TrimPrefix(s, "rclone://")
	log.Debugln("Opening " + path)

	r, w := io.Pipe()
	go func() {
		configfile.Install()
		fsrc := cmd.NewFsSrc([]string{path})
		err := operations.Cat(context.Background(), fsrc, w, 0, -1)
		if err != nil {
			log.Println(err)
		}
		_ = w.Close()
	}()

	return r, nil
}

func (f rcloneOpener) Evaluate(s string) float32 {
	// log.Debugf("Evaluating %s...\n", s)
	if strings.HasPrefix(s, "rclone://") {
		return 0.9
	}
	return 0
}
