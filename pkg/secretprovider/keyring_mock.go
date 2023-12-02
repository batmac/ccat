//go:build !keystore
// +build !keystore

package secretprovider

import (
	"errors"

	"braces.dev/errtrace"
	"github.com/batmac/ccat/pkg/log"
)

var (
	IsKeystoreAvailable = false
	ErrNotCompiled      = errors.New("keystore not compiled in")
)

func SetSecret(name, _ string) error {
	log.Printf("SetSecret(%s) called, but keystore is not compiled in", name)
	return errtrace.Wrap(ErrNotCompiled)
}

func getSecret(name string) (string, error) {
	log.Printf("GetSecret(%s) called, but keystore is not compiled in", name)
	return "", errtrace.Wrap(ErrNotCompiled)
}
