package mutators

import (
	"io"

	"github.com/ulikunitz/xz"
)

func init() {
	simpleRegister("unxz", unxz, withDescription("decompress xz data"))
	simpleRegister("xz", cxz, withDescription("compress xz data"), withExpectingBinary(true))
}

func unxz(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	d, err := xz.NewReader(in)
	if err != nil {
		return 0, err
	}
	n, err := io.Copy(out, d)
	return n, err
}

func cxz(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	h, err := xz.NewWriter(out)
	if err != nil {
		return 0, err
	}
	defer h.Close()
	return io.Copy(h, in)
}
