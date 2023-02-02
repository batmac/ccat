package mutators_test

import (
	"os/exec"
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

func Test_SingleExec(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"hello", "hello", "hello"},
		{"empty", "", ""},
		{"zero", "\x00", "\x00"},
	}

	var f string
	// search for a command that is available on both windows and linux
	if _, err := exec.LookPath("cat"); err == nil {
		f = "x:cat"
	} else if _, err := exec.LookPath("type"); err == nil {
		// on windows, "type" is probably available
		f = "x:type con"
	} else {
		t.Skip("no suitable command found")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); got != tt.encoded {
				t.Errorf("%s = %v, want %v", f, got, tt.encoded)
			}
		})
	}
}
