package mutators

import (
	"io"

	"github.com/batmac/ccat/pkg/log"
	"github.com/dsnet/compress/bzip2"
)

func init() {
	singleRegister("bzip2", cbzip2, withDescription("compress to bzip2 data (X:9 is compression level, 0-9)"),
		withCategory("compress"),
		withConfigBuilder(stdConfigUint64WithDefault(9)),
	)
	// singleRegister("unbzip2alt", bunzip2Alt, withDescription("decompress bzip2 data (alt)"))
}

func cbzip2(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	lvl := int(config.(uint64))
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
