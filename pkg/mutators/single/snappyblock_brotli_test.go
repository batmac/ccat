package mutators_test

import (
	"os"
	"strings"
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

// The fixtures come from the reference encoders, not from ccat:
//   - test.txt.br: Google's brotli (Python binding), quality 11
//   - test.txt.snappyblock: python-snappy, raw block
//   - test.txt.x50.xerial: kafka-python (xerial_compatible=True, 32KiB
//     blocks), 50 copies of test.txt, so it spans 4 chunks
func TestDecompressReferenceFixtures(t *testing.T) {
	content, err := os.ReadFile("testdata/compression/test.txt")
	if err != nil {
		t.Fatal(err)
	}
	x50 := strings.Repeat(string(content), 50)

	tests := []struct {
		mutator, fixture, want string
	}{
		{"unbrotli", "test.txt.br", string(content)},
		{"unsnapblock", "test.txt.snappyblock", string(content)},
		{"uns2block", "test.txt.snappyblock", string(content)},
		{"unxerial", "test.txt.x50.xerial", x50},
		// unxerial falls back to a raw block when there is no xerial header
		{"unxerial", "test.txt.snappyblock", string(content)},
	}
	for _, tt := range tests {
		t.Run(tt.mutator+" "+tt.fixture, func(t *testing.T) {
			in, err := os.ReadFile("testdata/compression/" + tt.fixture)
			if err != nil {
				t.Fatal(err)
			}
			if got := mutators.Run(tt.mutator, string(in)); got != tt.want {
				t.Errorf("mismatch (len %d, want %d)", len(got), len(tt.want))
			}
		})
	}
}

func TestXerialHeader(t *testing.T) {
	got := mutators.Run("xerial", "hello")
	// magic, then version 1 and min compatible version 1 (big endian)
	const header = "\x82SNAPPY\x00" + "\x00\x00\x00\x01" + "\x00\x00\x00\x01"
	if !strings.HasPrefix(got, header) {
		t.Errorf("missing xerial header: %q", got)
	}
}

func TestBlockAndXerialEdgeCases(t *testing.T) {
	for _, alg := range []string{"snapblock", "s2block", "xerial", "brotli"} {
		t.Run(alg+" empty", func(t *testing.T) {
			if got := mutators.Run("un"+alg, mutators.Run(alg, "")); got != "" {
				t.Errorf("round trip of empty input = %q", got)
			}
		})
	}
	input := bigInput() // several xerial chunks and brotli windows
	for _, alg := range []string{"snapblock", "s2block", "xerial", "brotli:1"} {
		t.Run(alg+" big", func(t *testing.T) {
			name, _, _ := strings.Cut(alg, ":")
			if got := mutators.Run("un"+name, mutators.Run(alg, input)); got != input {
				t.Errorf("round trip mismatch (len %d, want %d)", len(got), len(input))
			}
		})
	}
}
