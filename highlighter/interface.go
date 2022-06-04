package highlighter

import (
	"ccat/log"
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
}

func Go(w io.WriteCloser, r io.ReadCloser, o Options) error {
	c := new(Chroma)
	go func() {
		err := c.HighLight(w, r, o)
		if err != nil {
			log.Printf(" chroma highlighter returned an err: %v", err)
		}
		w.Close()
	}()
	return nil
}
