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
