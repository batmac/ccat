package stringutils_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/stringutils"
)

func TestHumanSize(t *testing.T) {
	tests := []struct {
		name  string
		input int64
		want  string
	}{
		{"zero", 0, "0.00"},
		{"one", 1, "1.00"},
		{"less than 1K", 999, "999.00"},
		{"one K", 1000, "1.00K"},
		{"one and half K", 1500, "1.50K"},
		{"one M", 1000000, "1.00M"},
		{"one G", 1000000000, "1.00G"},
		{"large number", 1234567, "1.23M"},
		// Add more cases as needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: Adapt the type if testing types other than int64
			if got := stringutils.HumanSize(tt.input); got != tt.want {
				t.Errorf("HumanSize(%d) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFromHumanSize(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{"zero", "0", 0, false},
		{"zero with unit", "0K", 0, false},
		{"hundred", "100", 100, false},
		{"one K", "1K", 1000, false},
		{"one and half K", "1.5K", 1500, false}, // Based on docker/go-units behavior
		{"one M", "1M", 1000000, false},
		{"one G", "1G", 1000000000, false},
		{"space K", "1 K", 1000, false},
		{"invalid unit", "1.5X", 0, true}, // Invalid unit
		{"just unit", "K", 0, true},       // Invalid format
		{"empty", "", 0, true},            // Invalid format
		{"non-numeric", "abc", 0, true},   // Invalid format
		// Add more cases as needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: Adapt the type T if testing types other than int64
			got, err := stringutils.FromHumanSize[int64](tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromHumanSize(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("FromHumanSize(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}
