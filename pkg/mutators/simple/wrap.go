package mutators

import (
	"bytes"
	"fmt"
	"io"

	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/wordwrap"
	uwrap "github.com/muesli/reflow/wrap"
)

const (
	WrapMaxChars = 80
	IndentChars  = 4
)

func init() {
	simplestRegister("wrap", wordWrap, withDescription(
		fmt.Sprintf("word-wrap the text to %d chars maximum", WrapMaxChars)))
	simplestRegister("wrapU", unconditionalyWrap, withDescription(
		fmt.Sprintf("unconditionally wrap the text to %d chars maximum", WrapMaxChars)))

	simplestRegister("indent", simpleIndent, withDescription(
		fmt.Sprintf("indent the text with %d chars", IndentChars)))
}

func wordWrap(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	ww := wordwrap.NewWriter(WrapMaxChars)
	if _, err := io.Copy(ww, r); err != nil { // streamable?
		return 0, err
	}
	ww.Close()
	return io.Copy(w, bytes.NewReader(ww.Bytes()))
}

func unconditionalyWrap(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	ww := uwrap.NewWriter(WrapMaxChars)
	if _, err := io.Copy(ww, r); err != nil { // streamable?
		return 0, err
	}

	return io.Copy(w, bytes.NewReader(ww.Bytes()))
}

func simpleIndent(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	f := indent.NewWriter(IndentChars, nil)
	if _, err := io.Copy(f, r); err != nil { // streamable?
		return 0, err
	}
	// f.Close()
	return io.Copy(w, bytes.NewReader(f.Bytes()))
}
