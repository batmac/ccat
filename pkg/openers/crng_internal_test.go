//go:build !fileonly
// +build !fileonly

package openers

import (
	"testing"
)

func Test_crngOpener_Evaluate(t *testing.T) {
	tests := []struct {
		name string
		args string
		want float32
	}{
		{name: "ko", args: "false://", want: 0},
		{name: "ok", args: "crng://", want: 0.9},
		{name: "ok with good limit", args: "crng://10", want: 0.9},
		{name: "ok with good limit 2", args: "crng://10k", want: 0.9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := crngOpener{}
			if got := f.Evaluate(tt.args); got != tt.want {
				t.Errorf("crngOpener.Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
