package mutators_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

func Test_discard(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"empty", "", ""},
		{"simple", "{}", ""},
		{"abc", "abc", ""},
		{"3x nl", "\n\n\n", ""},
	}

	f := "d"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); got != tt.encoded {
				t.Errorf("%s = %v, want %v", f, got, tt.encoded)
			}
		})
	}
}

func Test_wc(t *testing.T) {
	tests := []struct {
		name, input, defaultOption, b, r, w, l string
	}{
		{"empty", "", "0\n", "0\n", "0\n", "0\n", "0\n"},
		{"byte zero", "\x00", "1\n", "1\n", "1\n", "1\n", "1\n"},
		{"simple", "hello world !", "13\n", "13\n", "13\n", "3\n", "1\n"},
		{"C string", "hello\x00", "6\n", "6\n", "6\n", "1\n", "1\n"},
		{"invalid C string", "hello\x00world", "11\n", "11\n", "11\n", "1\n", "1\n"},
		{"abc ô", "abc ô", "6\n", "6\n", "5\n", "2\n", "1\n"},
		{"💋\n\n👍🏿", "💋\n\n👍🏿", "14\n", "14\n", "5\n", "2\n", "3\n"},
		{"3x nl", "\n\n\n", "3\n", "3\n", "3\n", "0\n", "3\n"},
		/*   		{"abc", "abc", ""},
		{"3x nl", "\n\n\n", ""}, */
	}

	f := "wc"
	for _, tt := range tests {
		for _, mode := range []string{"defaultOption", "b", "r", "w", "l"} {
			expected := map[string]string{"defaultOption": tt.defaultOption, "b": tt.b, "r": tt.r, "w": tt.w, "l": tt.l}
			t.Run(tt.name, func(t *testing.T) {
				mutatorName := f + ":" + mode
				if mode == "defaultOption" {
					mutatorName = f
				}
				if got := mutators.Run(mutatorName, tt.input); got != expected[mode] {
					t.Errorf("%s: %s = %v, for mode '%s' this is expected to be %v", tt.input, f, got, mode, expected[mode])
				}
			})
		}
	}
}
