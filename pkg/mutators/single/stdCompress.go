package mutators

import (
	"compress/bzip2"
	"compress/gzip"
	"compress/zlib"
	"io"
	"log"
)

func init() {
	singleNoConfRegister("ungzip", ungzip, withDescription("decompress gzip data"),
		withCategory("decompress"),
	)
	singleNoConfRegister("unbzip2", bunzip2, withDescription("decompress bzip2 data"),
		withCategory("decompress"),
	)
	singleNoConfRegister("unzlib", unzlib, withDescription("decompress zlib data"),
		withCategory("decompress"),
	)

	singleNoConfRegister("gzip", cgzip, withDescription("compress to gzip data"),
		withCategory("compress"),
	)
	singleNoConfRegister("zlib", czlib, withDescription("compress to zlib data"),
		withCategory("compress"),
	)
}

func ungzip(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		log.Fatal(err)
	}
	defer zr.Close()
	//#nosec
	return io.Copy(w, zr)
}

func bunzip2(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	bzr := bzip2.NewReader(r)
	if bzr == nil {
		log.Fatal("bzip2 decompressor failed to init")
	}
	//#nosec
	return io.Copy(w, bzr)
}

func unzlib(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	z, err := zlib.NewReader(r)
	if err != nil {
		log.Fatal(err)
	}
	defer z.Close()
	//#nosec
	return io.Copy(w, z)
}

func cgzip(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	zw, err := gzip.NewWriterLevel(w, gzip.DefaultCompression)
	if err != nil {
		log.Fatal(err)
	}
	defer zw.Close()
	//#nosec
	return io.Copy(zw, r)
}

func czlib(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	zw, err := zlib.NewWriterLevel(w, zlib.DefaultCompression)
	if err != nil {
		log.Fatal(err)
	}
	defer zw.Close()
	//#nosec
	return io.Copy(zw, r)
}
