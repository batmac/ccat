package mutators

import (
	"io"

	"github.com/batmac/ccat/pkg/log"
	"github.com/klauspost/compress/flate"
)

// raw DEFLATE (RFC 1951): no gzip/zlib header nor checksum.
func init() {
	singleRegister("undeflate", undeflate, withDescription("decompress raw deflate data (RFC 1951, no header)"),
		withCategory("decompress"),
		withAliases("inflate"),
	)
	singleRegister("deflate", cdeflate, withDescription("compress to raw deflate data (RFC 1951, no header) (X:6 is compression level, 0-9)"),
		withCategory("compress"),
		withConfigBuilder(stdConfigUint64WithDefault(^uint64(0))),
	)
}

func undeflate(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	fr := flate.NewReader(r)
	defer fr.Close()
	//#nosec
	return io.Copy(w, fr)
}

func cdeflate(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	lvl := cfgInt(config)
	log.Debugf("compression level: %d", lvl)
	fw, err := flate.NewWriter(w, lvl)
	if err != nil {
		return 0, err
	}
	n, err := io.Copy(fw, r)
	if err != nil {
		fw.Close()
		return n, err
	}
	return n, fw.Close()
}
