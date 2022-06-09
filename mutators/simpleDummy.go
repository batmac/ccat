package mutators

import (
	"io"
)

func init() {
	simpleRegister("dummy", "a simple fifo", "", dummy)
}

func dummy(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	return io.Copy(w, r)
}
