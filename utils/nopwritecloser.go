package utils

import "strings"

type NopWriteCloser struct {
	strings.Builder
}

func (NopWriteCloser) Close() error {
	return nil
}
