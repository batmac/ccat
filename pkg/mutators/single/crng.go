package mutators

import (
	"crypto/rand"
	"errors"
	"io"
)

func init() {
	singleRegister("crng", crng, withDescription("generate endless crypto/rand data "),
		withCategory("random"),
		withExpectingBinary(true),
	)
}

func crng(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	defer r.Close()
	defer w.Close()

	n, err := io.Copy(w, rand.Reader)
	if errors.Is(err, io.ErrClosedPipe) {
		err = nil
	}
	return n, err
}
