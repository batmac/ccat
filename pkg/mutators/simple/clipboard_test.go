package mutators_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
	"github.com/batmac/ccat/pkg/utils"
)

func Test_simpleClipboard(t *testing.T) {
	t.Run("donotpanicplease", func(t *testing.T) {
		f := "cb"
		if got := mutators.Run(f, " "); utils.DeleteSpaces(got) != "" {
			t.Errorf("%s = '%v', want '%v'", f, got, "")
		}
	})
}
