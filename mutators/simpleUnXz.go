package mutators

import (
	"io"

	"github.com/ulikunitz/xz"
)

func init() {
	simpleRegister("unxz", "decompresses xz data", decompress)
}

func decompress(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	d, err := xz.NewReader(in)
	if err != nil {
		return 0, err
	}
	n, err := io.Copy(out, d)
	return n, err
}
