package mutators

import (
	"braces.dev/errtrace"
	"io"
)

func init() {
	singleRegister("dummy", dummy,
		withDescription("a simple fifo"),
		withAliases("dum", "dumm"),
	)
}

func dummy(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	return errtrace.Wrap2(io.Copy(w, r)) // streamable
}
