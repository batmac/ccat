package mutators

import (
	"fmt"
	"io"

	"github.com/andybalholm/brotli"
	"github.com/batmac/ccat/pkg/log"
)

func init() {
	singleRegister("unbrotli", unbrotli, withDescription("decompress brotli data"),
		withCategory("decompress"),
	)
	singleRegister("brotli", cbrotli, withDescription("compress to brotli data (X:6 is compression level, 0-11)"),
		withCategory("compress"),
		withConfigBuilder(stdConfigUint64WithDefault(brotli.DefaultCompression)),
	)
}

func unbrotli(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	//#nosec
	return io.Copy(out, brotli.NewReader(in))
}

func cbrotli(out io.WriteCloser, in io.ReadCloser, config any) (int64, error) {
	lvl := cfgInt(config)
	// the library does not validate the level
	if lvl < brotli.BestSpeed || lvl > brotli.BestCompression {
		return 0, fmt.Errorf("brotli compression level must be %d-%d, got %d", brotli.BestSpeed, brotli.BestCompression, lvl)
	}
	log.Debugf("compression level: %d", lvl)
	w := brotli.NewWriterLevel(out, lvl)
	n, err := io.Copy(w, in)
	if err != nil {
		w.Close()
		return n, err
	}
	return n, w.Close()
}
