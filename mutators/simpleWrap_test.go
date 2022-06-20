package mutators_test

import (
	"testing"

	"github.com/batmac/ccat/mutators"
)

func Test_simpleWWrap(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{
			"hello",
			"hello",
			"hello",
		},
		{"empty", "", ""},

		{
			"lipsum",
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus. Suspendisse lectus tortor," +
				" dignissim sit amet, adipiscing nec, ultricies sed, dolor.",
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus.\n" +
				"Suspendisse lectus tortor, dignissim sit amet, adipiscing nec, ultricies sed,\n" +
				"dolor.",
		},
	}

	f := "wrap"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); got != tt.encoded {
				t.Errorf("%s = %v, want %v", f, got, tt.encoded)
			}
		})
	}
}

func Test_simpleUWrap(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{
			"hello",
			"hello",
			"hello",
		},
		{"empty", "", ""},

		{
			"lipsum",
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus. Suspendisse lectus tortor," +
				" dignissim sit amet, adipiscing nec, ultricies sed, dolor.",
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus. Suspendi\n" +
				"sse lectus tortor, dignissim sit amet, adipiscing nec, ultricies sed, dolor.",
		},
	}

	f := "wrapU"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); got != tt.encoded {
				t.Errorf("%s = %v, want %v", f, got, tt.encoded)
			}
		})
	}
}

func Test_simpleIndent(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{
			"hello",
			"hello",
			"    hello",
		},
		{"empty", "", ""},

		{
			"lipsum",
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus. Suspendisse lectus tortor,\n" +
				" dignissim sit amet, adipiscing nec, ultricies sed, dolor.",
			"    Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus. Suspendisse lectus tortor,\n" +
				"     dignissim sit amet, adipiscing nec, ultricies sed, dolor.",
		},
	}

	f := "indent"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); got != tt.encoded {
				t.Errorf("%s = %v, want %v", f, got, tt.encoded)
			}
		})
	}
}
