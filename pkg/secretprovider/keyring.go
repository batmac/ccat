//go:build keystore
// +build keystore

package secretprovider

import (
	"sync"

	"github.com/batmac/ccat/pkg/log"
	"github.com/tmc/keyring"
)

type (
	secretSrc struct {
		name   string
		envVar string
	}
	secretCollection struct {
		list map[string]secretSrc
		mu   sync.Mutex
	}
)

var IsKeystoreAvailable = true

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
