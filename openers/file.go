package openers

import (
	"ccat/lockable"
	"ccat/log"
	"io"
	"os"
	"strings"
)

var fileOpenerName = "file"
var fileOpenerDescription = "open local files"

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
	var err error
	if s != "-" {
		from, err = lockable.FileOpen(s, lock)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		//defer lockable.FileClose(from.(*os.File), false)
	}

	return from, nil
}

func (f fileOpener) Evaluate(s string) float32 {
	path := parsePath(s)
	//log.Debugf("returning %v", path)
	if path == "-" {
		return 1.0
	}
	_, err := os.Stat(path)
	if err == nil {
		return 0.99
	}
	return 0.0
}

func parsePath(s string) string {
	if strings.HasPrefix(s, "file://") {
		before, after, found := strings.Cut(s, "file://")
		//log.Debugf("before=%v after=%v, found=%v", before, after, found)
		if found && before == "" {
			return after
		}
	}
	return s
}
