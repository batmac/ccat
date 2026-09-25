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
	singleRegister("unbzip2", bunzip2, withDescription("decompress bzip2 data (in parallel, one worker per CPU)"),
		withCategory("decompress"),
	)
	singleRegister("punbzip2", punzip2, withDescription("parallel decompress bzip2 data (X:0 is concurrency, 0 is auto)"),
		withCategory("decompress"),
		withConfigBuilder(stdConfigUint64WithDefault(0)),
	)
}

func cbzip2(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	lvl := cfgInt(config)
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

func bunzip2(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	return pbunzip2(w, r, 0)
}

func punzip2(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	return pbunzip2(w, r, cfgInt(config))
}

// pbunzip2 decodes bzip2 blocks concurrently; about 5x faster than
// compress/bzip2 on 4 cores, and it handles concatenated streams too.
func pbunzip2(w io.Writer, r io.Reader, concurrency int) (int64, error) {
	if concurrency <= 0 {
		concurrency = runtime.NumCPU()
	}
	return io.Copy(w, pbzip2.NewReader(context.Background(), r, pbzip2.DecompressionOptions(pbzip2.BZConcurrency(concurrency))))
}
