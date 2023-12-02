package secretprovider

import (
	"errors"
	"os"

	"braces.dev/errtrace"
	"github.com/batmac/ccat/pkg/log"
	"github.com/tmc/keyring"
)

var (
	ErrNotFound = errors.New("secret not found")
	ServiceName = "ccat"
)

func GetSecret(name, envVar string) (string, error) {
	if v := os.Getenv(envVar); v != "" {
		log.Debugf("Using env var %s for secret '%s'", envVar, name)
		return v, nil
	}
	if IsKeystoreAvailable {
		s, err := getSecret(name)
		if errors.Is(err, keyring.ErrNotFound) {
			log.Debugf("Secret '%s' not found in keyring", name)
			return "", errtrace.Wrap(ErrNotFound)
		}
		return s, errtrace.Wrap(err)
	}
	return "", errtrace.Wrap(ErrNotFound)
}
