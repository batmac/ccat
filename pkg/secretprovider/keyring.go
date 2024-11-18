//go:build keystore
// +build keystore

package secretprovider

import (
	"github.com/batmac/ccat/pkg/log"
	"github.com/tmc/keyring"
)

var IsKeystoreAvailable = true
var ErrNotFound = keyring.ErrnotFound

func SetSecret(name, secret string) error {
	log.Debugf("Setting secret '%s'", name)
	err := keyring.Set(secretName(name), name, secret)
	if err != nil {
		log.Printf("Error setting secret '%s' in keyring: %v", name, err)
	}
	return err
}

func getSecret(name string) (string, error) {
	log.Debugf("Getting secret '%s'", name)
	s, err := keyring.Get(secretName(name), name)
	if err != nil {
		log.Debugf("Error getting secret '%s' from keyring: %v", name, err)
	}
	return s, err
}

func secretName(name string) string {
	return ServiceName + " " + name
}
