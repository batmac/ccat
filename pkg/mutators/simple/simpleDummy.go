package mutators

import (
	"io"
)

func init() {
	simpleRegister("dummy", dummy, withDescription("a simple fifo"))
}

func dummy(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	return io.Copy(w, r) // streamable
}
