package main

import (
	"os"

	"github.com/batmac/ccat/pkg/lockable"
)

func setupStdout(lock bool) error {
	if lock {
		err := lockable.Flock(os.Stdout)
		if err != nil {
			return err
		}
	}
	return nil
}
