package miniclaude_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/miniclaude"
)

func TestNewSimpleSamplingParameters(t *testing.T) {
	type args struct {
		prompt string
		model  string
	}
	tests := []struct {
		name string
		args args
	}{
		{"donotpanicplease", args{"hi", miniclaude.ModelClaudeLatest}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := miniclaude.NewSimpleSamplingParameters(tt.args.prompt, tt.args.model); got == nil {
				t.Errorf("NewSimpleSamplingParameters() = nil")
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"donotpanicplease"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := miniclaude.New(); got == nil {
				t.Errorf("New() = nil")
			}
		})
	}
}
