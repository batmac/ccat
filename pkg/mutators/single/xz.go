package mutators

import (
	"io"

	"braces.dev/errtrace"
	"github.com/ulikunitz/xz"
	"github.com/ulikunitz/xz/lzma"
)

func init() {
	singleRegister("unxz", unxz, withDescription("decompress xz data"),
		withCategory("decompress"),
	)
	singleRegister("xz", cxz, withDescription("compress to xz data"),
		withCategory("compress"),
	)

	singleRegister("unlzma", unlzma, withDescription("decompress lzma data"),
		withCategory("decompress"),
	)
	singleRegister("lzma", clzma, withDescription("compress to lzma data"),
		withCategory("compress"),
	)

	singleRegister("unlzma2", unlzma2, withDescription("decompress lzma2 data"),
		withCategory("decompress"),
	)
	singleRegister("lzma2", clzma2, withDescription("compress to lzma2 data"),
		withCategory("compress"),
	)
}

func unxz(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	d, err := xz.NewReader(in)
	if err != nil {
		return 0, errtrace.Wrap(err)
	}
	n, err := io.Copy(out, d)
	return n, errtrace.Wrap(err)
}

func cxz(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	h, err := xz.NewWriter(out)
	if err != nil {
		return 0, errtrace.Wrap(err)
	}
	defer h.Close()
	return errtrace.Wrap2(io.Copy(h, in))
}

func unlzma(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	d, err := lzma.NewReader(in)
	if err != nil {
		return 0, errtrace.Wrap(err)
	}
	n, err := io.Copy(out, d)
	return n, errtrace.Wrap(err)
}

func clzma(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	h, err := lzma.NewWriter(out)
	if err != nil {
		return 0, errtrace.Wrap(err)
	}
	defer h.Close()
	return errtrace.Wrap2(io.Copy(h, in))
}

func unlzma2(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	d, err := lzma.NewReader2(in)
	if err != nil {
		return 0, errtrace.Wrap(err)
	}
	n, err := io.Copy(out, d)
	return n, errtrace.Wrap(err)
}

func clzma2(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	h, err := lzma.NewWriter2(out)
	if err != nil {
		return 0, errtrace.Wrap(err)
	}
	defer h.Close()
	return errtrace.Wrap2(io.Copy(h, in))
}
