package mutators

import (
	"ccat/log"
	"io"

	"github.com/klauspost/compress/s2"
)

func init() {
	simpleRegister("uns2", "decompress s2 data", "", uns2)
	simpleRegister("s2", "compress s2 data", "", cs2)

	simpleRegister("unsnap", "decompress snappy data", "", uns2)
	simpleRegister("snap", "compress snappy data", "", csnappy)

}

func uns2(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	d := s2.NewReader(in)
	if d == nil {
		log.Fatal("s2 decompressor failed to init")
	}
	n, err := io.Copy(out, d)
	return n, err
}

func cs2(dst io.WriteCloser, src io.ReadCloser) (int64, error) {
	return _cs2(dst, src)

}

func csnappy(dst io.WriteCloser, src io.ReadCloser) (int64, error) {
	return _cs2(dst, src, s2.WriterSnappyCompat())
}

func _cs2(dst io.WriteCloser, src io.ReadCloser, opts ...s2.WriterOption) (int64, error) {
	enc := s2.NewWriter(dst, opts...)
	n, err := io.Copy(enc, src)
	if err != nil {
		enc.Close()
		return 0, err
	}
	// Blocks until compression is done.
	enc.Close()
	return n, nil
}
