//go:build !windows
// +build !windows

package lockable

import (
	"os"
	"syscall"
)

func Flock(file *os.File) error {
	return syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
}

func Unflock(file *os.File) error {
	return syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
}
