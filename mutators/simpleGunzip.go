package mutators

import (
	"compress/bzip2"
	"compress/gzip"
	"io"
	"log"
)

func init() {
	simpleRegister("gunzip", "decompresses gzip data", gunzip)

	simpleRegister("bzip2", "decompresses bzip2 data", bunzip2)
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
