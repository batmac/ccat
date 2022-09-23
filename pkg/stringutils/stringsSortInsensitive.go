package stringutils

import (
	"sort"
	"unicode"
	"unicode/utf8"
)

// https://stackoverflow.com/questions/51997276/how-one-can-do-case-insensitive-sorting-using-sort-strings-in-golang
// https://gist.github.com/Deleplace/8a9224f8b2838a9c1d112acde5bdd3c8

// lessCaseInsensitive compares s, t without allocating
func lessCaseInsensitive(s, t string) bool {
	for {
		if len(t) == 0 {
			return false
		}
		if len(s) == 0 {
			return true
		}
		c, sizec := utf8.DecodeRuneInString(s)
		d, sized := utf8.DecodeRuneInString(t)

		lowerc := unicode.ToLower(c)
		lowerd := unicode.ToLower(d)

		if lowerc < lowerd {
			return true
		}

		if lowerc > lowerd {
			return false
		}

		s = s[sizec:]
		t = t[sized:]
	}
}

func SortStringsCaseInsensitive(data []string) {
	sort.Slice(data, func(i, j int) bool { return lessCaseInsensitive(data[i], data[j]) })
}
