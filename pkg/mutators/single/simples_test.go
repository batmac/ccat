package mutators_test

import (
	"fmt"
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
		name, input string
		expectedCounts map[string]int64 // mode -> expected count
	}{
		{"empty", "", map[string]int64{"b": 0, "r": 0, "w": 0, "l": 0}},
		{"byte zero", "\x00", map[string]int64{"b": 1, "r": 1, "w": 1, "l": 1}},
		{"simple", "hello world !", map[string]int64{"b": 13, "r": 13, "w": 3, "l": 1}},
		{"C string", "hello\x00", map[string]int64{"b": 6, "r": 6, "w": 1, "l": 1}},
		{"invalid C string", "hello\x00world", map[string]int64{"b": 11, "r": 11, "w": 1, "l": 1}},
		{"abc ô", "abc ô", map[string]int64{"b": 6, "r": 5, "w": 2, "l": 1}},
		{"💋\n\n👍🏿", "💋\n\n👍🏿", map[string]int64{"b": 14, "r": 5, "w": 2, "l": 3}},
		{"3x nl", "\n\n\n", map[string]int64{"b": 3, "r": 3, "w": 0, "l": 3}},
	}

	f := "wc"
	for _, tt := range tests {
		for _, mode := range []string{"b", "r", "w", "l"} {
			expectedCount := tt.expectedCounts[mode]
			t.Run(tt.name, func(t *testing.T) {
				mutatorName := f + ":" + mode
				got := mutators.Run(mutatorName, tt.input)
				
				// The output could be either "13\n" or "13.00\n" depending on terminal detection
				// We'll validate that the numeric value is correct regardless of formatting
				var actualCount int64
				if n, err := fmt.Sscanf(got, "%d\n", &actualCount); n == 1 && err == nil {
					// Integer format like "13\n"
					if actualCount != expectedCount {
						t.Errorf("%s: %s = %d, for mode '%s' expected count %d", tt.input, f, actualCount, mode, expectedCount)
					}
				} else {
					// Try float format like "13.00\n"
					var floatCount float64
					if n, err := fmt.Sscanf(got, "%f\n", &floatCount); n == 1 && err == nil {
						actualCount = int64(floatCount)
						if actualCount != expectedCount {
							t.Errorf("%s: %s = %d, for mode '%s' expected count %d", tt.input, f, actualCount, mode, expectedCount)
						}
					} else {
						t.Errorf("%s: %s = %q, for mode '%s' could not parse as number", tt.input, f, got, mode)
					}
				}
			})
		}
		
		// Test default mode (should be same as 'b' mode)
		t.Run(tt.name, func(t *testing.T) {
			expectedCount := tt.expectedCounts["b"]
			got := mutators.Run(f, tt.input)
			
			var actualCount int64
			if n, err := fmt.Sscanf(got, "%d\n", &actualCount); n == 1 && err == nil {
				if actualCount != expectedCount {
					t.Errorf("%s: %s = %d, for default mode expected count %d", tt.input, f, actualCount, expectedCount)
				}
			} else {
				// Try float format like "13.00\n"
				var floatCount float64
				if n, err := fmt.Sscanf(got, "%f\n", &floatCount); n == 1 && err == nil {
					actualCount = int64(floatCount)
					if actualCount != expectedCount {
						t.Errorf("%s: %s = %d, for default mode expected count %d", tt.input, f, actualCount, expectedCount)
					}
				} else {
					t.Errorf("%s: %s = %q, for default mode could not parse as number", tt.input, f, got)
				}
			}
		})
	}
}
