package color_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/color"
)

func TestANSI_Next(t *testing.T) {
	testNext(t, new(color.ANSI))
}

func TestANSIbg_Next(t *testing.T) {
	testNext(t, new(color.ANSIbg))
}

func Test256_Next(t *testing.T) {
	testNext(t, new(color.C256))
}

func testNext(t *testing.T, c color.Color) {
	t.Helper()
	t.Run("donotpanicplease", func(t *testing.T) {
		source := "hi"
		for i := 0; i < 100; i++ {
			c = c.Next()
			s := c.Sprint(source)
			if len(s) <= len(source) {
				t.FailNow()
			}
		}
	})
}
