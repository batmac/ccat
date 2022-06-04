package main

import (
	"ccat/lockable"
	"ccat/log"
	"os"
)

func setupStdout(lock bool) {
	if lock {
		err := lockable.Flock(os.Stdout)
		if err != nil {
			log.Fatal(err)
		}
	}
}
