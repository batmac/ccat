//go:build !nohl
// +build !nohl

package highlighter

import (
	"io"

	"github.com/batmac/ccat/log"
)

type Options struct {
	FileName      string
	StyleHint     string
	LexerHint     string
	FormatterHint string
}

type highlighter interface {
	highLight(w io.WriteCloser, r io.ReadCloser, o Options) error
	help() string
}

func Go(w io.WriteCloser, r io.ReadCloser, o Options) error {
	go func() {
		c := new(Chroma)
		err := c.highLight(w, r, o)
		if err != nil {
			log.Printf(" chroma highlighter returned an err: %v", err)
		}
		w.Close()
	}()
	return nil
}

func Help() string {
	c := new(Chroma)
	return c.help()
}
