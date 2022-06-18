package mutators

import (
	"github.com/dsnet/compress/bzip2"

	"io"
	"log"
)

func init() {
	simpleRegister("bzip2", cbzip2, withDescription("compress bzip2 data"), withExpectingBinary(true))
	// simpleRegister("unbzip2alt", bunzip2Alt, withDescription("decompress bzip2 data (alt)"))
}

func cbzip2(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	zw, err := bzip2.NewWriter(w, &bzip2.WriterConfig{Level: bzip2.BestCompression})
	if err != nil {
		log.Fatal(err)
	}
	defer zw.Close()
	return io.Copy(zw, r)
}

/* func bunzip2Alt(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	bzr, err := bzip2.NewReader(r, nil)
	if err != nil {
		log.Fatal(err)
	}
	return io.Copy(w, bzr)
} */
