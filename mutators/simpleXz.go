package mutators

import (
	"io"

	"github.com/ulikunitz/xz"
	"github.com/ulikunitz/xz/lzma"
)

func init() {
	simpleRegister("unxz", unxz, withDescription("decompress xz data"))
	simpleRegister("xz", cxz, withDescription("compress xz data"), withExpectingBinary(true))

	simpleRegister("unlzma", unlzma, withDescription("decompress lzma data"))
	simpleRegister("lzma", clzma, withDescription("compress lzma data"), withExpectingBinary(true))
	simpleRegister("unlzma2", unlzma2, withDescription("decompress lzma2 data"))
	simpleRegister("lzma2", clzma2, withDescription("compress lzma2 data"), withExpectingBinary(true))

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

func unlzma(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	d, err := lzma.NewReader(in)
	if err != nil {
		return 0, err
	}
	n, err := io.Copy(out, d)
	return n, err
}

func clzma(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	h, err := lzma.NewWriter(out)
	if err != nil {
		return 0, err
	}
	defer h.Close()
	return io.Copy(h, in)
}

func unlzma2(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	d, err := lzma.NewReader2(in)
	if err != nil {
		return 0, err
	}
	n, err := io.Copy(out, d)
	return n, err
}

func clzma2(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	h, err := lzma.NewWriter2(out)
	if err != nil {
		return 0, err
	}
	defer h.Close()
	return io.Copy(h, in)
}
