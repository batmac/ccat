package mutators_test

import (
	"strings"
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

func Test_simpleJ5J(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"empty", "", ""},
		{"simple", "{}", "{}"},
		{"indented", `{hi:'hi',}`, "{\n  \"hi\": \"hi\"\n}"},
		{"comment", `
{
	hi:'hi'
	//hi
}
`, "{\n  \"hi\": \"hi\"\n}"},
		{"https://json5.org example", `{
	// comments
	unquoted: 'and you can quote me on that',
	simpleQuotes: 'I can use "double quotes" here',
	lineBreaks: "Look, Mom! \
  No \\n's!",
	hexadecimal: 0xdecaf,
	leadingDecimalPoint: .8675309, andTrailing: 8675309.,
	positiveSign: +1,
	trailingComma: 'in objects', andIn: ['arrays',],
	"backwardsCompatible": "with JSON",
  }`, `{
	"andIn": [
	  "arrays"
	],
	"andTrailing": 8675309,
	"backwardsCompatible": "with JSON",
	"hexadecimal": 912559,
	"leadingDecimalPoint": 0.8675309,
	"lineBreaks": "Look, Mom!   No \\n's!",
	"positiveSign": 1,
	"simpleQuotes": "I can use \"double quotes\" here",
	"trailingComma": "in objects",
	"unquoted": "and you can quote me on that"
  }
  `},
		{"indented2", `{"hi": 1.}`, "{\n  \"hi\": 1\n}"},
		{"indenteds1", `{"hi": +1}`, "{\n  \"hi\": 1\n}"},
		{"indentedhex", `{"hi": 0x1}`, "{\n  \"hi\": 1\n}"},
		{"indented3", "   { \n \"hi\" :    1 \n    }", "{\n  \"hi\": 1\n}"},
		{"indented4", "  \n\n { \n\n \n  \n \"hi\"  \n: \n   1 \n}", "{\n  \"hi\": 1\n}"},
	}

	f := "j5j"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); mutators.Run("jcs", got) != mutators.Run("jcs", tt.encoded) {
				t.Errorf("%s = %v, want %v", f, strings.TrimSuffix(got, "\n"), tt.encoded)
			}
		})
	}
}
