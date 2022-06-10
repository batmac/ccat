package mutators

import (
	"compress/bzip2"
	"compress/gzip"
	"compress/zlib"
	"io"
	"log"
)

func init() {
	simpleRegister("gunzip", gunzip, withDescription("decompress gzip data"))
	simpleRegister("bunzip2", bunzip2, withDescription("decompress bzip2 data"))
	simpleRegister("unzlib", unzlib, withDescription("decompress zlib data"))
}

func gunzip(w io.WriteCloser, r io.ReadCloser) (int64, error) {
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
	return io.Copy(w, r)

}
