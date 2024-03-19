package miniclaude_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/miniclaude"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"donotpanicplease"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := miniclaude.NewMessagesRequest(); got == nil {
				t.Errorf("New() = nil")
			}
		})
	}
}
