package mutators_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

// float64s builds a byte stream of count doubles, as ccat would receive it.
func float64s(count int, f func(i int) float64) string {
	var b bytes.Buffer
	for i := range count {
		_ = binary.Write(&b, binary.LittleEndian, f(i))
	}
	return b.String()
}

func TestShuffleRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		size  int
		input string
	}{
		{"empty", 8, ""},
		{"single element", 8, "abcdefgh"},
		{"shorter than one element", 8, "abc"},
		{"not a multiple of the element size", 8, "abcdefghij"},
		{"element size 1 is a no-op", 1, "hello world"},
		{"element size 3", 3, "aaabbbcccddd"},
		{"element size 4", 4, float64s(64, func(i int) float64 { return float64(i) })},
		{"doubles", 8, float64s(1000, func(i int) float64 { return 100 + math.Sin(float64(i)/10) })},
		// larger than one 64KiB block, to exercise the block loop
		{"multi-block", 8, float64s(20000, func(i int) float64 { return float64(i) * 1.5 })},
		{"multi-block, partial last block", 8, float64s(20000, func(i int) float64 { return float64(i) }) + "xyz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shuffled := mutators.Run(fmt.Sprintf("shuffle:%d", tt.size), tt.input)
			if len(shuffled) != len(tt.input) {
				t.Fatalf("shuffle changed the length: got %d, want %d", len(shuffled), len(tt.input))
			}
			got := mutators.Run(fmt.Sprintf("unshuffle:%d", tt.size), shuffled)
			if got != tt.input {
				t.Errorf("round-trip mismatch: got %d bytes, want %d", len(got), len(tt.input))
			}
		})
	}
}

func TestShuffleGroupsBytePositions(t *testing.T) {
	// four 4-byte elements: the transform must emit all the first bytes, then
	// all the second ones, and so on
	input := "AbcdAbcdAbcdAbcd"
	want := "AAAAbbbbccccdddd"
	if got := mutators.Run("shuffle:4", input); got != want {
		t.Errorf("shuffle:4 = %q, want %q", got, want)
	}
}

func TestShuffleHelpsCompression(t *testing.T) {
	// the point of the filter: a compressor that cannot do anything with
	// interleaved doubles does well once they are grouped by byte position
	input := float64s(4000, func(i int) float64 { return 1000 + math.Sin(float64(i)/50) })

	plain := len(mutators.Run("zstd", input))
	shuffled := len(mutators.Run("zstd", mutators.Run("shuffle:8", input)))

	t.Logf("%d bytes -> zstd %d, shuffle+zstd %d", len(input), plain, shuffled)
	if shuffled >= plain {
		t.Errorf("shuffle did not help: shuffle+zstd = %d, zstd alone = %d", shuffled, plain)
	}
}

func TestShuffleDefaultsToEightBytes(t *testing.T) {
	input := strings.Repeat("Abcdefgh", 4)
	if mutators.Run("shuffle", input) != mutators.Run("shuffle:8", input) {
		t.Error("shuffle should default to an element size of 8")
	}
}
