package mutators

import (
	"io"

	"github.com/batmac/ccat/pkg/log"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zlib"
)

func init() {
	// unpgzip is an alias: pgzip's read-ahead reader is not faster than
	// klauspost/gzip, pgzip only pays off when compressing
	singleRegister("ungzip", ungzip, withDescription("decompress gzip data"),
		withCategory("decompress"),
		withAliases("unpgzip"),
	)
	singleRegister("unzlib", unzlib, withDescription("decompress zlib data"),
		withCategory("decompress"),
	)

	singleRegister("gzip", cgzip, withDescription("compress to gzip data (X:6 is compression level, 0-9)"),
		withCategory("compress"),
		withConfigBuilder(stdConfigUint64WithDefault(^uint64(0))),
	)
	singleRegister("zlib", czlib, withDescription("compress to zlib data (X:6 is compression level, 0-9)"),
		withCategory("compress"),
		withConfigBuilder(stdConfigUint64WithDefault(^uint64(0))),
	)
}

func ungzip(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		log.Fatal(err)
	}
	defer zr.Close()
	//#nosec
	return io.Copy(w, zr)
}

func unzlib(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	z, err := zlib.NewReader(r)
	if err != nil {
		log.Fatal(err)
	}
	defer z.Close()
	//#nosec
	return io.Copy(w, z)
}

func cgzip(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	lvl := cfgInt(config)
	// log.Printf("compression level: %d", lvl)
	zw, err := gzip.NewWriterLevel(w, lvl)
	if err != nil {
		log.Fatal(err)
	}
	defer zw.Close()
	//#nosec
	return io.Copy(zw, r)
}

func czlib(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	lvl := cfgInt(config)
	log.Debugf("compression level: %d", lvl)
	zw, err := zlib.NewWriterLevel(w, lvl)
	if err != nil {
		log.Fatal(err)
	}
	defer zw.Close()
	//#nosec
	return io.Copy(zw, r)
}
