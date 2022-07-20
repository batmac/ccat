package openers

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/batmac/ccat/pkg/lockable"
	"github.com/batmac/ccat/pkg/log"
)

var (
	fileOpenerName        = "file"
	fileOpenerDescription = "open local files"
)

type fileOpener struct {
	name, description string
}

func init() {
	register(&fileOpener{
		name:        fileOpenerName,
		description: fileOpenerDescription,
	})
}

func (f fileOpener) Name() string {
	return f.name
}

func (f fileOpener) Description() string {
	return f.description
}

func (f fileOpener) Open(s string, lock bool) (io.ReadCloser, error) {
	s = parsePath(s)
	var from io.ReadCloser = os.Stdin
	// var err error
	if s != "-" {
		fileInfo, err := os.Stat(s)
		if err != nil {
			return nil, err
		}
		if fileInfo.IsDir() {
			return nil, errors.New("Is a directory")
		}
		from, err = lockable.FileOpen(s, lock)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		// defer lockable.FileClose(from.(*os.File), false)
	}

	return from, nil
}

func (f fileOpener) Evaluate(s string) float32 {
	path := parsePath(s)
	// log.Debugf("returning %v", path)
	if path == "-" {
		return 1.0
	}
	if _, err := os.Stat(path); err == nil {
		return 0.99
	}
	return 0.0
}

func parsePath(s string) string {
	if strings.HasPrefix(s, "file://") {
		after := s[7:]
		// log.Debugf("before=%v after=%v, found=%v", before, after, found)
		return after
	}
	return s
}
