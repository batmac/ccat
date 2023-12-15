package mutators_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

func Test_mistral(t *testing.T) {
	// only test that we do not panic
	t.Setenv("MISTRAL_API_KEY", "CI")

	f := "mistral:100:fakemodel"
	t.Run("donotpanicplease", func(t *testing.T) {
		if got := mutators.Run(f, "hi"); got != "CI" {
			t.Errorf("%s = %v, want %v", f, got, "CI")
		}
	})
}
