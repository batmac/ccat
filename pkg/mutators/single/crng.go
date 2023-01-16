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
		withConfigBuilder(stdConfigHumanSizeAsInt64WithDefault(int64(0))),
	)
}

func crng(w io.WriteCloser, r io.ReadCloser, conf any) (int64, error) {
	N := conf.(int64)
	defer r.Close()
	defer w.Close()

	var n int64
	var err error
	if N > 0 {
		n, err = io.CopyN(w, rand.Reader, N)
	} else {
		n, err = io.Copy(w, rand.Reader)
	}
	if errors.Is(err, io.ErrClosedPipe) {
		err = nil
	}
	return n, err
}
