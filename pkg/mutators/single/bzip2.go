package mutators

import (
	"context"
	"io"
	"runtime"

	"github.com/batmac/ccat/pkg/log"
	"github.com/cosnicolaou/pbzip2"
	"github.com/dsnet/compress/bzip2"
)

func init() {
	singleRegister("bzip2", cbzip2, withDescription("compress to bzip2 data (X:9 is compression level, 0-9)"),
		withCategory("compress"),
		withConfigBuilder(stdConfigUint64WithDefault(9)),
	)
	singleRegister("punbzip2", punzip2, withDescription("parallel decompress bzip2 data (X:0 is concurrency, 0 is auto)"),
		withCategory("decompress"),
		withConfigBuilder(stdConfigUint64WithDefault(0)),
	)
}

func cbzip2(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	lvl := int(config.(uint64)) //nolint:gosec
	log.Debugf("compression level: %d", lvl)
	zw, err := bzip2.NewWriter(w, &bzip2.WriterConfig{Level: lvl})
	if err != nil {
		log.Fatal(err)
	}
	defer zw.Close()
	return io.Copy(zw, r) // streamable
}

/* func bunzip2Alt(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	bzr, err := bzip2.NewReader(r, nil)
	if err != nil {
		log.Fatal(err)
	}
	return io.Copy(w, bzr)
} */

func punzip2(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	concurrency := int(config.(uint64)) //nolint:gosec
	if concurrency == 0 {
		concurrency = runtime.NumCPU()
	}
	return io.Copy(w, pbzip2.NewReader(context.Background(), r, pbzip2.DecompressionOptions(pbzip2.BZConcurrency(concurrency))))
}
