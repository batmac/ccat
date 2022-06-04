package highlighter

import (
	"errors"
	"io"
)

//import "github.com/rubenv/pygmentize"
type Pygment struct {
}

func (h Pygment) HihLight(w io.WriteCloser, r io.ReadCloser, o Options) error {
	return errors.New("Not implemented")
}
