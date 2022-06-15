package mutators_test

import (
	"testing"

	"github.com/batmac/ccat/mutators"
)

func Test_simpleFilterUTF8(t *testing.T) {
	var tests = []struct {
		name, decoded, encoded string
	}{
		{"empty", "", ""},
		{"simple", "a\xe9b", "ab"},
		{"Valid ASCII", "a", "a"},
		{"Valid ASCII", "\x5A\x6F\xEB", "Zo"},
		{"Valid 2 Octet Sequence", "X\xc3\xb1X", "XñX"},
		{"Invalid 2 Octet Sequence", "X\xc3\x28X", "X(X"},
		{"Invalid Sequence Identifier", "X\xa0\xa1X", "XX"},
		{"Valid 3 Octet Sequence", "X\xe2\x82\xa1X", "X₡X"},
		{"Invalid 3 Octet Sequence (in 2nd Octet)", "X\xe2\x28\xa1X", "X(X"},
		{"Invalid 3 Octet Sequence (in 3rd Octet)", "X\xe2\x82\x28X", "X(X"},
		{"Valid 4 Octet Sequence", "X\xf0\x90\x8c\xbcX", "X𐌼X"},
		{"Invalid 4 Octet Sequence (in 2nd Octet)", "X\xf0\x28\x8c\xbcX", "X(X"},
		{"Invalid 4 Octet Sequence (in 3rd Octet)", "X\xf0\x90\x28\xbcX", "X(X"},
		{"Invalid 4 Octet Sequence (in 4th Octet)", "X\xf0\x28\x8c\x28X", "X((X"},
		{"Valid 5 Octet Sequence (but not Unicode!)", "X\xf8\xa1\xa1\xa1\xa1X", "XX"},
		{"Valid 6 Octet Sequence (but not Unicode!)", "X\xfc\xa1\xa1\xa1\xa1\xa1X", "XX"},
	}

	f := "filterUTF8"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); got != tt.encoded {
				t.Errorf("%s = %s, want %s", f, got, tt.encoded)
			}
		})
	}
}
