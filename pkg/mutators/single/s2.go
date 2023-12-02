package mutators

import (
	"io"

	"github.com/batmac/ccat/pkg/log"

	"braces.dev/errtrace"
	"github.com/klauspost/compress/s2"
)

func init() {
	singleRegister("uns2", uns2, withDescription("decompress s2 data"),
		withCategory("decompress"),
	)
	singleRegister("s2", cs2, withDescription("compress to s2 data"),
		withCategory("compress"),
	)

	singleRegister("unsnap", uns2, withDescription("decompress snappy data"),
		withCategory("decompress"),
	)
	singleRegister("snap", csnappy, withDescription("compress to snappy data"),
		withCategory("compress"),
	)
}

func uns2(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	d := s2.NewReader(in)
	if d == nil {
		log.Fatal("s2 decompressor failed to init")
	}
	n, err := io.Copy(out, d)
	return n, errtrace.Wrap(err)
}

func cs2(dst io.WriteCloser, src io.ReadCloser, _ any) (int64, error) {
	return errtrace.Wrap2(_cs2(dst, src))
}

func csnappy(dst io.WriteCloser, src io.ReadCloser, _ any) (int64, error) {
	return errtrace.Wrap2(_cs2(dst, src, s2.WriterSnappyCompat()))
}

func _cs2(dst io.WriteCloser, src io.ReadCloser, opts ...s2.WriterOption) (int64, error) {
	enc := s2.NewWriter(dst, opts...)
	n, err := io.Copy(enc, src)
	if err != nil {
		enc.Close()
		return 0, errtrace.Wrap(err)
	}
	// Blocks until compression is done.
	enc.Close()
	return n, nil
}
