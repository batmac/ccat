package mutators_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

func Test_claude(t *testing.T) {
	// only test that we do not panic
	t.Setenv("ANTHROPIC_API_KEY", "CI")

	f := "claude:100:fakemodel"
	t.Run("donotpanicplease", func(t *testing.T) {
		if got := mutators.Run(f, "hi"); got != "fake" {
			t.Errorf("%s = %v, want %v", f, got, "fake")
		}
	})
}
