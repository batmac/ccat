package mutators

import (
	"io"
	"runtime"

	gzip "github.com/klauspost/pgzip"

	"github.com/batmac/ccat/pkg/log"
)

func init() {
	singleRegister("unpgzip", unpgzip, withDescription("decompress with pgzip"),
		withCategory("decompress"),
	)

	singleRegister("pgzip", cpgzip, withDescription("compress with pgzip  (X:6 is compression level, 0-9, blockSize, blocks)"),
		withCategory("compress"),
		withConfigBuilder(stdConfigInts(0, 3)),
	)
}

func unpgzip(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		log.Fatal(err)
	}
	defer zr.Close()

	//#nosec
	return io.Copy(w, zr)
}

func cpgzip(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	args := config.([]int)

	lvl := 6
	if len(args) > 0 {
		lvl = args[0]
	}

	log.Debugf("compression level: %d", lvl)
	zw, err := gzip.NewWriterLevel(w, lvl)
	if err != nil {
		log.Fatal(err)
	}
	defer zw.Close()

	switch len(args) {
	case 2:
		log.Debugf("setting block size: %d", args[1])
		if err := zw.SetConcurrency(args[1], runtime.GOMAXPROCS(0)); err != nil {
			log.Fatal(err)
		}
	case 3:
		log.Debugf("setting block size: %d, blocks: %d", args[1], args[2])
		if err := zw.SetConcurrency(args[1], args[2]); err != nil {
			log.Fatal(err)
		}
	}

	//#nosec
	return io.Copy(zw, r)
}
