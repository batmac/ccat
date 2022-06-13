package mutators

import (
	"bytes"
	"fmt"
	"io"

	"github.com/muesli/reflow/wordwrap"
	uwrap "github.com/muesli/reflow/wrap"
)

const WRAP_MAX_CHARS = 80

func init() {
	simpleRegister("wrap", wordWrap, withDescription(fmt.Sprintf("word-wrap the text to %d chars maximum", WRAP_MAX_CHARS)))
	simpleRegister("wrapU", unconditionalyWrap, withDescription(fmt.Sprintf("unconditionaly wrap the text to %d chars maximum", WRAP_MAX_CHARS)))
}

func wordWrap(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	ww := wordwrap.NewWriter(WRAP_MAX_CHARS)
	io.Copy(ww, r)
	ww.Close()
	return io.Copy(w, bytes.NewReader(ww.Bytes()))
}
func unconditionalyWrap(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	ww := uwrap.NewWriter(WRAP_MAX_CHARS)
	io.Copy(ww, r)

	return io.Copy(w, bytes.NewReader(ww.Bytes()))
}
