package mutators_test

import (
	"encoding/binary"
	"strings"
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

func Test_bzip3Header(t *testing.T) {
	tests := []struct {
		name      string
		blockSize uint32
	}{
		{"bzip3", 16 << 20},
		{"bzip3:1MB", 1_000_000},
		{"bzip3:66560", 65 << 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mutators.Run(tt.name, "hello")
			if !strings.HasPrefix(got, "BZ3v1") || len(got) < 9 {
				t.Fatalf("%s: missing BZ3v1 header: %q", tt.name, got)
			}
			if bs := binary.LittleEndian.Uint32([]byte(got[5:9])); bs != tt.blockSize {
				t.Errorf("%s: block size = %d, want %d", tt.name, bs, tt.blockSize)
			}
		})
	}
}

func Test_bzip3MultiBlockRoundTrip(t *testing.T) {
	// 3+ blocks at the minimum block size
	input := strings.Repeat("The quick brown fox jumps over the lazy dog. ", 5000)
	if got := mutators.Run("unbzip3", mutators.Run("bzip3:66560", input)); got != input {
		t.Errorf("round trip mismatch: len = %d, want %d", len(got), len(input))
	}
}

func Test_bzip3Empty(t *testing.T) {
	if got := mutators.Run("unbzip3", mutators.Run("bzip3", "")); got != "" {
		t.Errorf("round trip of empty input = %q", got)
	}
}

func Test_bzip3InvalidBlockSize(t *testing.T) {
	for _, arg := range []string{"1k", "1GB", "notasize"} {
		if _, err := mutators.New("bzip3:" + arg); err == nil {
			t.Errorf("bzip3:%s: expected an error", arg)
		}
	}
}
