package lockable

import (
	"os"
	"path/filepath"

	"braces.dev/errtrace"
	"github.com/batmac/ccat/pkg/log"
)

// open and optionally flock a file
func FileOpen(path string, lock bool) (*os.File, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	log.Debugln(" lockable: opened ", file.Name())
	if lock {
		err = Flock(file)
		if err != nil {
			file.Close()
			return nil, errtrace.Wrap(err)
		}
		log.Println(" lockable: locked ", file.Name())
	}
	return file, nil
}

// optionally unflock and close a file
func FileClose(file *os.File, unlock bool) {
	if unlock {
		log.Debugln(" lockable: unlock ", file.Name())
		err := Unflock(file)
		if err != nil {
			log.Println(err)
		}
	}
	log.Debugln(" lockable: close ", file.Name())
	if err := file.Close(); err != nil {
		log.Println(err)
	}
}
