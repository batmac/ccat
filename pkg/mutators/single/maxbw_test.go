package mutators_test

import (
	"testing"
	"time"

	"github.com/batmac/ccat/pkg/mutators"
)

func Test_maxbw(t *testing.T) {
	tests := []struct {
		name, input, output string
	}{
		{"empty", "", ""},
		{"salut", "salut", "salut"},
	}
	f := "maxbw:1G"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime := time.Now()
			if got := mutators.Run(f, tt.input); tt.output != got {
				t.Errorf("%s = '%v', want %v", f, got, tt.output)
			}
			diff := time.Since(startTime)
			t.Logf("diff:%v", diff)
		})
	}
}
