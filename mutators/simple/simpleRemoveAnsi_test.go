package mutators_test

import (
	"testing"

	"github.com/batmac/ccat/mutators"
)

func Test_simpleRemoveANSI(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"hello", "hello", "hello"},
		{"empty", "", ""},
		{"zero", "\x00", "\x00"},
		{"line", "\x1b[7m\x1b[0mchar\x1b[93;41m\nnl", "char\nnl"},
	}

	f := "removeANSI"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); got != tt.encoded {
				t.Errorf("%s = %v, want %v", f, got, tt.encoded)
			}
		})
	}
}
