package term_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/term"
)

func TestOsc52(t *testing.T) {
	something := []byte("This is a test string from ccat testing (TestOsc52)")
	t.Run("donotpanicplease", func(t *testing.T) {
		term.Osc52(something)
	})
}
