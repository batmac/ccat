//go:build !fileonly
// +build !fileonly

package openers

import (
	"testing"
)

func Test_prngOpener_Evaluate(t *testing.T) {
	tests := []struct {
		name string
		args string
		want float32
	}{
		{name: "ko", args: "false://", want: 0},
		{name: "ok", args: "prng://", want: 0.9},
		{name: "ok with good limit", args: "prng://10", want: 0.9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := prngOpener{}
			if got := f.Evaluate(tt.args); got != tt.want {
				t.Errorf("prngOpener.Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
