package utils_test

import (
	"strings"
	"testing"

	"github.com/batmac/ccat/pkg/utils"
)

func TestNopStringWriteCloser_Close(t *testing.T) {
	t.Run("donotpanicplease", func(t *testing.T) {
		n := utils.NopStringWriteCloser{strings.Builder{}}
		if err := n.Close(); err != nil {
			t.Errorf("NopStringWriteCloser.Close() error = %v", err)
		}
	})
}
