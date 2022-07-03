package term_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/term"
)

func TestOsc52(t *testing.T) {
	empty := []byte{0x0}
	t.Run("donotpanicplease", func(t *testing.T) {
		term.Osc52(empty)
	})
}
