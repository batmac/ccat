package stringutils

import (
	"strings"
	"unicode"
)

// https://stackoverflow.com/questions/32081808/strip-all-whitespace-from-a-string
func DeleteSpaces(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}
