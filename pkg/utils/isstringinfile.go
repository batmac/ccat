package utils

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

func IsStringInFile(s, path string) bool {
	if s == "" {
		panic("empty string")
	}
	// use only with small files as we read it fully
	d, err := os.ReadFile(filepath.Clean(path))
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	if err != nil {
		return false
	}
	return bytes.Contains(d, []byte(s))
}
