package mutators_test

import (
	"strings"
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

func Test_simpleJsonIndent(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"empty", "", ""},
		{"simple", "{}", "{}"},
		{"indented", `{"hi":"hi"}`, "{\n  \"hi\": \"hi\"\n}"},
		{"indented2", `{"hi": 1}`, "{\n  \"hi\": 1\n}"},
		{"indented3", "   { \n \"hi\" :    1 \n    }", "{\n  \"hi\": 1\n}"},
		{"indented4", "  \n\n { \n\n \n  \n \"hi\"  \n: \n   1 \n}", "{\n  \"hi\": 1\n}"},
	}

	f := "j"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); strings.TrimSuffix(got, "\n") != tt.encoded {
				t.Errorf("%s = %v, want %v", f, strings.TrimSuffix(got, "\n"), tt.encoded)
			}
		})
	}
}
