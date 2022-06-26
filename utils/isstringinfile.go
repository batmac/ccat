package utils

import (
	"bytes"
	"io/ioutil"
	"os"
)

func IsStringInFile(s, path string) bool {
	if len(s) <= 0 {
		panic("empty string")
	}
	// use only with small files as we read it fully
	d, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		panic("file does not exist")
	}
	if err != nil {
		panic(err)
	}
	return bytes.Contains(d, []byte(s))
}
