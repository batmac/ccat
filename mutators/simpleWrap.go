package mutators

import (
	"bytes"
	"fmt"
	"io"

	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/wordwrap"
	uwrap "github.com/muesli/reflow/wrap"
)

const WRAP_MAX_CHARS = 80
const INDENT_CHARS = 4

func init() {
	simpleRegister("wrap", wordWrap, withDescription(
		fmt.Sprintf("word-wrap the text to %d chars maximum", WRAP_MAX_CHARS)))
	simpleRegister("wrapU", unconditionalyWrap, withDescription(
		fmt.Sprintf("unconditionally wrap the text to %d chars maximum", WRAP_MAX_CHARS)))

	simpleRegister("indent", simpleIndent, withDescription(
		fmt.Sprintf("indent the text with %d chars", INDENT_CHARS)))

}

func wordWrap(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	ww := wordwrap.NewWriter(WRAP_MAX_CHARS)
	_, err := io.Copy(ww, r)
	if err != nil {
		return 0, err
	}
	ww.Close()
	return io.Copy(w, bytes.NewReader(ww.Bytes()))
}
func unconditionalyWrap(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	ww := uwrap.NewWriter(WRAP_MAX_CHARS)
	_, err := io.Copy(ww, r)
	if err != nil {
		return 0, err
	}

	return io.Copy(w, bytes.NewReader(ww.Bytes()))
}

func simpleIndent(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	f := indent.NewWriter(INDENT_CHARS, nil)
	_, err := io.Copy(f, r)
	if err != nil {
		return 0, err
	}
	//f.Close()
	return io.Copy(w, bytes.NewReader(f.Bytes()))
}
