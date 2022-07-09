package mutators_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

func Test_simpleTranslate(t *testing.T) {
	f := "translate"
	t.Run("donotpanicplease", func(t *testing.T) {
		if got := mutators.Run(f, ""); got != "" {
			t.Errorf("%s = %v, want %v", f, got, "")
		}
	})
}
