package mutators_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/mutators"
)

func TestNewGeneric(t *testing.T) {
	t.Run("donotpanicplease", func(t *testing.T) {
		if got := mutators.NewGeneric(log.Default()); got.Logger != log.Default() {
			t.Errorf("NewGeneric(), want log.Default()")
		}
	})

}
