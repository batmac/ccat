package lockable

import (
	"os"

	"github.com/batmac/ccat/log"
)

// open and optionally flock a file
func FileOpen(path string, lock bool) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	log.Debugln(" lockable: opened ", file.Name())
	if lock {
		err = Flock(file)
		if err != nil {
			file.Close()
			return nil, err
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
