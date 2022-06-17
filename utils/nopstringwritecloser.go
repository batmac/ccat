package utils

import (
	"strings"
)

type NopStringWriteCloser struct {
	strings.Builder
}

func (NopStringWriteCloser) Close() error {
	return nil
}
