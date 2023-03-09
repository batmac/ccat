package mutators_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

func Test_simpleDummy(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"hello", "hello", "hello"},
		{"empty", "", ""},
		{"zero", "\x00", "\x00"},
	}

	fl := []string{"dummy", "dum"}
	for _, f := range fl {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := mutators.Run(f, tt.decoded); got != tt.encoded {
					t.Errorf("%s = %v, want %v", f, got, tt.encoded)
				}
			})
		}
	}
}
