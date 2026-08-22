package stringutils

import "slices"

func IsStringInSlice(a string, list []string) bool {
	return slices.Contains(list, a)
}
