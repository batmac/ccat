package mutators

import (
	"io"

	"github.com/klauspost/compress/zstd"
)

func init() {
	simpleRegister("unzstd", "decompress zstd data", unzstd)
	simpleRegister("zstd", "compress zstd data", czstd)
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

// Create a reader that caches decompressors.
// For this operation type we supply a nil Reader.

// Decompress a buffer. We don't supply a destination buffer,
// so it will be allocated by the decoder.
func Decompress(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	var decoder, _ = zstd.NewReader(in, zstd.WithDecoderConcurrency(0))
	return io.Copy(out, decoder)
}
