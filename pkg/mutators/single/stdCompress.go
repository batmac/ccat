package mutators

import (
	"compress/bzip2"
	"compress/gzip"
	"compress/zlib"
	"io"

	"braces.dev/errtrace"
	"github.com/batmac/ccat/pkg/log"
)

func init() {
	singleRegister("ungzip", ungzip, withDescription("decompress gzip data"),
		withCategory("decompress"),
	)
	singleRegister("unbzip2", bunzip2, withDescription("decompress bzip2 data"),
		withCategory("decompress"),
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
	return errtrace.Wrap2(io.Copy(w, zr))
}

func bunzip2(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	bzr := bzip2.NewReader(r)
	if bzr == nil {
		log.Fatal("bzip2 decompressor failed to init")
	}
	//#nosec
	return errtrace.Wrap2(io.Copy(w, bzr))
}

func unzlib(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	z, err := zlib.NewReader(r)
	if err != nil {
		log.Fatal(err)
	}
	defer z.Close()
	//#nosec
	return errtrace.Wrap2(io.Copy(w, z))
}

func cgzip(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	lvl := int(config.(uint64))
	// log.Printf("compression level: %d", lvl)
	zw, err := gzip.NewWriterLevel(w, lvl)
	if err != nil {
		log.Fatal(err)
	}
	defer zw.Close()
	//#nosec
	return errtrace.Wrap2(io.Copy(zw, r))
}

func czlib(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	lvl := int(config.(uint64))
	log.Debugf("compression level: %d", lvl)
	zw, err := zlib.NewWriterLevel(w, lvl)
	if err != nil {
		log.Fatal(err)
	}
	defer zw.Close()
	//#nosec
	return errtrace.Wrap2(io.Copy(zw, r))
}
