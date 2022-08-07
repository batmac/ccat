//go:build !nohl
// +build !nohl

package highlighter

import (
	"io"
	"strings"

	"github.com/batmac/ccat/pkg/log"
)

type highlighter interface {
	highLight(w io.WriteCloser, r io.ReadCloser, o Options) error
	help() string
}

func Go(w io.WriteCloser, r io.ReadCloser, o Options) error {
	go func() {
		var c highlighter = new(Chroma)
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

func Run(input string, o *Options) string {
	in := io.NopCloser(strings.NewReader(input))
	r, w := io.Pipe()
	if o == nil {
		o = NewOptions("", "", "", "")
	}
	if err := Go(w, in, *o); err != nil {
		log.Printf("error while highlighting: %v", err)
	}
	reply, err := io.ReadAll(r)
	if err != nil {
		log.Printf("failed to read the highlighted string %v", err)
		return input
	}
	return string(reply)
}
