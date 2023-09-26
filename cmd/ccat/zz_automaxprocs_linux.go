package main

import (
	"github.com/batmac/ccat/pkg/log"
	"go.uber.org/automaxprocs/maxprocs"
)

// https://go.dev/ref/spec#Package_initialization
func init() {
	_, err := maxprocs.Set(maxprocs.Logger(log.Debugf))
	if err != nil {
		log.Debugf("failed to set GOMAXPROCS: %v", err)
	}
}
