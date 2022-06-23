package utils

import (
	"bytes"
	"io/ioutil"
)

func IsStringInFile(s, path string) bool {
	// use only with small files as we read it fully
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}
	return bytes.Contains(d, []byte(s))
}
