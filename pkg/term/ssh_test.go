package term_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/term"
)

func TestIsSsh(t *testing.T) {
	t.Run("donotpanicplease", func(t *testing.T) {
		_ = term.IsSsh()
	})
}
