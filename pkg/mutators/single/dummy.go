package mutators

import (
	"io"
)

func init() {
	singleRegister("dummy", dummy,
		withDescription("a simple fifo"),
		withAliases("dum", "dumm"),
	)
}

func dummy(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	return io.Copy(w, r) // streamable
}
