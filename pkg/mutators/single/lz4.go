package mutators

import (
	"io"

	"github.com/batmac/ccat/pkg/log"

	"github.com/pierrec/lz4/v4"
)

func init() {
	singleNoConfRegister("unlz4", unlz4, withDescription("decompress lz4 data"),
		withCategory("decompress"),
	)
	singleNoConfRegister("lz4", clz4, withDescription("compress to lz4 data"),
		withCategory("compress"),
	)
}

func unlz4(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	d := lz4.NewReader(in)
	if d == nil {
		log.Fatal("lz4 decompressor failed to init")
	}

	n, err := io.Copy(out, d)
	return n, err
}

func clz4(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	e := lz4.NewWriter(out)
	if e == nil {
		log.Fatal("lz4 compressor failed to init")
	}

	n, err := io.Copy(e, in)
	e.Close()
	return n, err
}
