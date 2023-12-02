package secretprovider_test

import (
	"errors"
	"testing"

	"github.com/batmac/ccat/pkg/secretprovider"
)

func TestGetSecret(t *testing.T) {
	t.Run("dontpanicplease", func(t *testing.T) {
		t.Setenv("CCAT_TEST_SECRET", "") // just to be sure
		got, err := secretprovider.GetSecret("nonexistingstuff", "CCAT_TEST_SECRET")
		if !errors.Is(err, secretprovider.ErrNotFound) {
			t.Errorf("GetSecret() failed: error = %v, got = %v", err, got)
			return
		}
	})

	t.Run("envar", func(t *testing.T) {
		t.Setenv("CCAT_TEST_SECRET", "test")
		got, err := secretprovider.GetSecret("nonexistingstuff", "CCAT_TEST_SECRET")
		if got != "test" || err != nil {
			t.Errorf("GetSecret() failed: error = %v, got = %v", err, got)
			return
		}
	})
}
