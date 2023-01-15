package mutators_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

func Test_singleWA(t *testing.T) {
	for _, f := range []string{"wa", "wasimple", "waspoken"} {
		t.Run("donotpanicplease", func(t *testing.T) {
			if got := mutators.Run(f, ""); got != "" {
				t.Errorf("%s = %v, want %v", f, got, "")
			}
		})
	}
}
