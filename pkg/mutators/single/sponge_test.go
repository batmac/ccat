package mutators_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

func Test_sponge(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"empty", "", ""},
		{"simple", "{}", "{}"},
		{"abc", "abc", "abc"},
		{"nl", "\n", "\n"},
		{"3x nl", "\n\n\n", "\n\n\n"},
	}

	f := "sponge"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); got != tt.encoded {
				t.Errorf("%s = %v, want %v", f, got, tt.encoded)
			}
		})
	}
}
