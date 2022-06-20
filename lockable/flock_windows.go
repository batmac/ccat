//go:build windows
// +build windows

package lockable

import (
	"os"

	"github.com/batmac/ccat/log"
)

func Flock(file *os.File) error {
	log.Println("flock is not supported on this platform, noop")
	return nil
}

func Unflock(file *os.File) error {
	return nil
}
