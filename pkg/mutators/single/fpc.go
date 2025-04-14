package mutators

import (
	"io"

	"github.com/spenczar/fpc"
)

func init() {
	singleRegister("fpc", wfpc, withDescription("compress 64bit floats"),
		withCategory("compress"),
	)
	singleRegister("unfpc", rfpc, withDescription("decompress 64bit floats"),
		withCategory("decompress"),
	)
}

func wfpc(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	f := fpc.NewWriter(out)

	// the buffer must obviously be a multiple of 8
	n, err := io.CopyBuffer(f, in, make([]byte, 8*1024))
	if err != nil {
		return n, err
	}
	err = f.Close()
	return n, err
}

func rfpc(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	return io.Copy(out, fpc.NewReader(in))
}
