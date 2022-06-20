package mutators

import (
	"compress/bzip2"
	"compress/gzip"
	"compress/zlib"
	"io"
	"log"
)

func init() {
	simpleRegister("ungzip", ungzip, withDescription("decompress gzip data"),
		withCategory("decompress"),
	)
	simpleRegister("unbzip2", bunzip2, withDescription("decompress bzip2 data"),
		withCategory("decompress"),
	)
	simpleRegister("unzlib", unzlib, withDescription("decompress zlib data"),
		withCategory("decompress"),
	)

	simpleRegister("gzip", cgzip, withDescription("compress to gzip data"),
		withCategory("compress"),
	)
	simpleRegister("zlib", czlib, withDescription("compress to zlib data"),
		withCategory("compress"),
	)
}

func ungzip(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		log.Fatal(err)
	}
	defer zr.Close()
	return io.Copy(w, zr)
}

func bunzip2(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	bzr := bzip2.NewReader(r)
	if bzr == nil {
		log.Fatal("bzip2 decompressor failed to init")
	}
	return io.Copy(w, bzr)
}

func unzlib(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	z, err := zlib.NewReader(r)
	if err != nil {
		log.Fatal(err)
	}
	defer z.Close()
	return io.Copy(w, z)

}

func cgzip(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	zw, err := gzip.NewWriterLevel(w, gzip.DefaultCompression)
	if err != nil {
		log.Fatal(err)
	}
	defer zw.Close()
	return io.Copy(zw, r)
}
func czlib(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	zw, err := zlib.NewWriterLevel(w, zlib.DefaultCompression)
	if err != nil {
		log.Fatal(err)
	}
	defer zw.Close()
	return io.Copy(zw, r)
}
