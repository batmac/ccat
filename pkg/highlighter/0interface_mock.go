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

func Go(w io.WriteCloser, r io.ReadCloser, _ Options) error {
	go func() {
		_, _ = io.Copy(w, r)
		w.Close()
	}()
	return nil
}

func Help() string {
	return "not supported (compiled with nohl)\n"
}
