package mutators_test

import (
	"fmt"
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

func Test_simpleLimit(t *testing.T) {
	limitedSize := 10
	tests := []struct {
		name, decoded, encoded string
	}{
		{"hello", "hello", "hello"},
		{"empty", "", ""},
		{"zero", "\x00", "\x00"},
		{
			"long",
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			"Lorem ipsu",
		},
	}

	f := "limit:" + fmt.Sprint(limitedSize)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); got != tt.encoded {
				t.Errorf("%s = %v, want %v", f, got, tt.encoded)
			}
		})
	}
}

func Test_simpleLimitWithSISize(t *testing.T) {
	limitedSizes := []string{
		"10k", "10m", "10g", "10t", "10p", "10b",
		"10kb", "10mb", "10gb", "10tb", "10pb",
		"10K", "10M", "10G", "10T", "10P", "10B",
		"10KB", "10MB", "10GB", "10TB", "10PB",
	}
	tests := []struct {
		name, decoded, encoded string
	}{
		{"hello", "hello", "hello"},
		{"empty", "", ""},
		{"zero", "\x00", "\x00"},
	}

	for _, limitedSize := range limitedSizes {
		f := "limit:" + limitedSize
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := mutators.Run(f, tt.decoded); got != tt.encoded {
					t.Errorf("%s = %v, want %v", f, got, tt.encoded)
				}
			})
		}
	}
}

func Test_simpleLimit0(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"hello", "hello", ""},
		{"empty", "", ""},
		{"zero", "\x00", ""},
		{
			"long",
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			"",
		},
	}

	f := "limit:0"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); got != tt.encoded {
				t.Errorf("%s = %v, want %v", f, got, tt.encoded)
			}
		})
	}
}
