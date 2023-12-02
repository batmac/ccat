package main

import (
	"os"

	"braces.dev/errtrace"
	"github.com/batmac/ccat/pkg/lockable"
)

func setupStdout(lock bool) error {
	if lock {
		err := lockable.Flock(os.Stdout)
		if err != nil {
			return errtrace.Wrap(err)
		}
	}
	return nil
}
