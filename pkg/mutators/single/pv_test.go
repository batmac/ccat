package mutators_test

import (
	"crypto/rand"
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

func Test_pv(t *testing.T) {
	// 100M of random data
	data := make([]byte, 100*1024*1024)
	_, _ = rand.Read(data)

	tests := []struct {
		name, input, expected string
	}{
		{"empty", "", ""},
		{"byte zero", "\x00", "\x00"},
		{"simple", "hello world !", "hello world !"},
		{"alphabet", "abcdef\x00ghijklmn\nopqr\rstuvwxyz", "abcdef\x00ghijklmn\nopqr\rstuvwxyz"},
		{"random", string(data), string(data)},
	}

	f := "pv:1"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.input); got != tt.expected {
				t.Errorf(" %s = %v, this is expected to be %v", tt.input, got, tt.expected)
			}
		})
	}
}
