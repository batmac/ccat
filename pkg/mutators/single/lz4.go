package mutators

import (
	"io"

	"github.com/batmac/ccat/pkg/log"

	"braces.dev/errtrace"
	"github.com/pierrec/lz4/v4"
)

func init() {
	singleRegister("unlz4", unlz4, withDescription("decompress lz4 data"),
		withCategory("decompress"),
	)
	singleRegister("lz4", clz4, withDescription("compress to lz4 data (X:0 is compression level, 0-9)"),
		withCategory("compress"),
		withConfigBuilder(stdConfigUint64WithDefault(uint64(lz4.Fast))),
	)
}

func unlz4(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	d := lz4.NewReader(in)
	if d == nil {
		log.Fatal("lz4 decompressor failed to init")
	}

	n, err := io.Copy(out, d)
	return n, errtrace.Wrap(err)
}

func clz4(out io.WriteCloser, in io.ReadCloser, config any) (int64, error) {
	c := config.(uint64)
	compressionLevel := lz4.Fast
	if c != 0 {
		compressionLevel = lz4.CompressionLevel(1 << (8 + lz4.CompressionLevel(c)))
	}
	log.Debugf("compression level: %s", compressionLevel)
	e := lz4.NewWriter(out)
	if e == nil {
		log.Fatal("compressor failed to init")
	}
	if err := e.Apply(lz4.CompressionLevelOption(compressionLevel)); err != nil {
		log.Fatal(err.Error())
	}

	n, err := io.Copy(e, in)
	e.Close()
	return n, errtrace.Wrap(err)
}
