package main

import (
	"ccat/log"
	"os"
	"syscall"
)

func setupStdout(lock bool) {
	if lock {
		err := syscall.Flock(int(os.Stdout.Fd()), syscall.LOCK_EX)
		if err != nil {
			log.Fatal(err)
		}
	}
}
