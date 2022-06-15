package utils

import (
	"strings"
	"unicode"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

/*
// https://stackoverflow.com/questions/7052693/how-to-get-the-name-of-a-function-in-go
func GetFunctionName(temp interface{}) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name()), ".")
	return strs[len(strs)-1]
}
*/
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
