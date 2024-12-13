package secretprovider

import (
	"errors"
	"os"

	"github.com/batmac/ccat/pkg/log"
)

var ServiceName = "ccat"

func GetSecret(name, envVar string) (string, error) {
	if v := os.Getenv(envVar); v != "" {
		log.Debugf("Using env var %s for secret '%s'", envVar, name)
		return v, nil
	}
	if IsKeystoreAvailable {
		s, err := getSecret(name)
		if errors.Is(err, ErrNotFound) {
			log.Debugf("Secret '%s' not found in keyring", name)
			return "", ErrNotFound
		}
		return s, err
	}
	return "", ErrNotFound
}
