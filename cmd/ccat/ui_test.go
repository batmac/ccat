package main

import (
	"io"
	"testing"
)

func Test_uiWrapProcessFile(t *testing.T) {
	t.Run("donotpanicplease", func(t *testing.T) {
		_ = uiWrapProcessFile(func(io.Writer, string) {})
	})
}
