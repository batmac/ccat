package mutators_test

import (
	"strings"
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

func Test_simpleCrc32(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"empty", "", "00000000"},
		{"abc", "abc", "352441c2"},
		{"simple", "{}", "a3a6bf43"},
		{"hello", "hello", "3610a686"},
		{"hello_world", "hello world", "0d4a1185"},
		{"123456789", "123456789", "cbf43926"},
		{"long_text", "The quick brown fox jumps over the lazy dog", "414fa339"},
		{"alphanumeric", "abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq", "171a3f5f"},
	}

	f := "crc32"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); strings.TrimSuffix(got, "\n") != tt.encoded {
				t.Errorf("%s = %v, want %v", f, strings.TrimSuffix(got, "\n"), tt.encoded)
			}
		})
	}
}