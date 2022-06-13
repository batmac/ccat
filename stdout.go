package main

import (
	"os"

	"github.com/batmac/ccat/lockable"
	"github.com/batmac/ccat/log"
)

func setupStdout(lock bool) {
	if lock {
		err := lockable.Flock(os.Stdout)
		if err != nil {
			log.Fatal(err)
		}
	}
}
