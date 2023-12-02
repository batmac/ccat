//go:build !windows
// +build !windows

package lockable

import (
	"braces.dev/errtrace"
	"os"
	"syscall"
)

func Flock(file *os.File) error {
	return errtrace.Wrap(syscall.Flock(int(file.Fd()), syscall.LOCK_EX))
}

func Unflock(file *os.File) error {
	return errtrace.Wrap(syscall.Flock(int(file.Fd()), syscall.LOCK_UN))
}
