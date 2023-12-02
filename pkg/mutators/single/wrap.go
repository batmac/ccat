package mutators

import (
	"bytes"
	"fmt"
	"io"

	"braces.dev/errtrace"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/wordwrap"
	uwrap "github.com/muesli/reflow/wrap"
)

const (
	WrapMaxChars = 80
	IndentChars  = 4
)

func init() {
	singleRegister("wrap", wordWrap, withDescription(
		fmt.Sprintf("word-wrap the text (to X:%d chars maximum)", WrapMaxChars)),
		withConfigBuilder(stdConfigUint64WithDefault(WrapMaxChars)))
	singleRegister("wrapU", unconditionalyWrap, withDescription(
		fmt.Sprintf("unconditionally wrap the text (to X:%d chars maximum)", WrapMaxChars)),
		withConfigBuilder(stdConfigUint64WithDefault(WrapMaxChars)))
	singleRegister("indent", singleIndent, withDescription(
		fmt.Sprintf("indent the text (with X:%d chars)", IndentChars)),
		withConfigBuilder(stdConfigUint64WithDefault(IndentChars)))
}

func wordWrap(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	WrapMaxChars := int(config.(uint64))
	ww := wordwrap.NewWriter(WrapMaxChars)
	if _, err := io.Copy(ww, r); err != nil { // streamable?
		return 0, errtrace.Wrap(err)
	}
	ww.Close()
	return errtrace.Wrap2(io.Copy(w, bytes.NewReader(ww.Bytes())))
}

func unconditionalyWrap(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	WrapMaxChars := int(config.(uint64))
	ww := uwrap.NewWriter(WrapMaxChars)
	if _, err := io.Copy(ww, r); err != nil { // streamable?
		return 0, errtrace.Wrap(err)
	}

	return errtrace.Wrap2(io.Copy(w, bytes.NewReader(ww.Bytes())))
}

func singleIndent(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	IndentChars := uint(config.(uint64))
	f := indent.NewWriter(IndentChars, nil)
	if _, err := io.Copy(f, r); err != nil { // streamable?
		return 0, errtrace.Wrap(err)
	}
	// f.Close()
	return errtrace.Wrap2(io.Copy(w, bytes.NewReader(f.Bytes())))
}
