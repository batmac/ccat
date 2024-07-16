package completion_test

import (
	_ "embed"
	"testing"

	"github.com/batmac/ccat/pkg/completion"
)

func TestPrint(t *testing.T) {
	type args struct {
		shell string
		opts  []string
	}
	tests := []struct {
		name string
		args args
	}{
		{"bash", args{"bash", []string{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			completion.Print(tt.args.shell, tt.args.opts)
		})
	}
}
