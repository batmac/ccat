package utils

import (
	"braces.dev/errtrace"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func FileChecksum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", errtrace.Wrap(err)
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", errtrace.Wrap(err)
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
