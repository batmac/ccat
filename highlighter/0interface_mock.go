//go:build nohl
// +build nohl

package highlighter

import (
	"io"
)

type Options struct {
	FileName      string
	StyleHint     string
	LexerHint     string
	FormatterHint string
}

type Highlighter interface {
	HighLight(w io.WriteCloser, r io.ReadCloser, o Options) error
	help() string
}

func Go(w io.WriteCloser, r io.ReadCloser, o Options) error {
	return nil
}

func Help() string {
	return "not supported (compiled with nohl)\n"
}
