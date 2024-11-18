//go:build !keystore
// +build !keystore

package secretprovider

import (
	"errors"

	"github.com/batmac/ccat/pkg/log"
)

var (
	IsKeystoreAvailable = false
	ErrNotCompiled      = errors.New("keystore not compiled in")
	ErrNotFound         = errors.New("Not Found")
)

func SetSecret(name, _ string) error {
	log.Printf("SetSecret(%s) called, but keystore is not compiled in", name)
	return ErrNotCompiled
}

func getSecret(name string) (string, error) {
	log.Printf("GetSecret(%s) called, but keystore is not compiled in", name)
	return "", ErrNotCompiled
}
