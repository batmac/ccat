package stringutils_test

import (
	"strings"
	"testing"

	"github.com/batmac/ccat/pkg/stringutils"
)

func TestNopStringWriteCloser_Close(t *testing.T) {
	t.Run("donotpanicplease", func(t *testing.T) {
		n := stringutils.NopStringWriteCloser{strings.Builder{}}
		if err := n.Close(); err != nil {
			t.Errorf("NopStringWriteCloser.Close() error = %v", err)
		}
	})
}
