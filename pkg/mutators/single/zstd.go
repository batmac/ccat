package mutators

import (
	"archive/zip"
	"io"

	"github.com/klauspost/compress/zstd"
)

func init() {
	registerAsZipDecompressor()

	singleNoConfRegister("unzstd", unzstd, withDescription("decompress zstd data"),
		withCategory("decompress"),
	)
	singleNoConfRegister("zstd", czstd, withDescription("compress to zstd data"),
		withCategory("compress"),
	)
}

func registerAsZipDecompressor() {
	decomp := zstd.ZipDecompressor()
	zip.RegisterDecompressor(zstd.ZipMethodWinZip, decomp)
	zip.RegisterDecompressor(zstd.ZipMethodPKWare, decomp)
}

func unzstd(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	d, err := zstd.NewReader(in)
	if err != nil {
		return 0, err
	}
	defer d.Close()

	n, err := io.Copy(out, d)
	return n, err
}

func czstd(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	e, err := zstd.NewWriter(out)
	if err != nil {
		return 0, err
	}

	n, err := io.Copy(e, in)
	e.Close()
	return n, err
}
