//go:build crappy
// +build crappy

package mutators_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/mutators"
)

func Test_simpleEasySeal(t *testing.T) {
	t.Run("donotpanicplease", func(t *testing.T) {
		f := "easyseal"
		if got := mutators.Run(f, ""); got == "" {
			t.Errorf("%s = '%v', want '%v'", f, got, "")
		}
	})
}

func Test_simpleEasyOpen(t *testing.T) {
	t.Run("donotpanicplease", func(t *testing.T) {
		f := "easyopen"
		log.SetContinueOnFatal()
		if got := mutators.Run(f, ""); got != "" {
			t.Errorf("%s = '%v', want '%v'", f, got, "")
		}
	})
}

func Test_simpleNaclE2E(t *testing.T) {
	t.Run("e2e", func(t *testing.T) {
		// len(32) key (64 chars hex)
		hexKey := "1234567890123456789012345678901234567890123456789012345678901234"
		t.Setenv("KEY", hexKey)

		seal := mutators.Run("easyseal", "hello")
		if seal == "" {
			t.Fatalf("mutators.Run() error : seal is empty")
		}

		opened := mutators.Run("easyopen", seal)

		if opened != "hello" {
			t.Errorf("mutators.Run() = '%v', want '%v'", opened, "hello")
		}
	})
}
