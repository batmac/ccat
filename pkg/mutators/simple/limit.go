package mutators

import (
	"io"
	"strconv"
)

func init() {
	simpleRegister("limit", limit, withDescription("a simple limiting fifo"))
}

func limit(w io.WriteCloser, r io.ReadCloser, args ...string) (int64, error) {
	bytes, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return 0, err
	}

	lr := io.LimitReader(r, bytes)

	return io.Copy(w, lr) // streamable
}
