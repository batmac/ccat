package mutators

import (
	"archive/zip"
	"io"

	"github.com/batmac/ccat/pkg/log"
	"github.com/klauspost/compress/zstd"
)

func init() {
	registerAsZipDecompressor()

	singleRegister("unzstd", unzstd, withDescription("decompress zstd data"),
		withCategory("decompress"),
	)
	singleRegister("zstd", czstd, withDescription("compress to zstd data (X:4 is compression level, 1-22)"),
		withCategory("compress"),
		withConfigBuilder(stdConfigUint64WithDefault(uint64(zstd.SpeedDefault))),
	)
}

func registerAsZipDecompressor() {
	decomp := zstd.ZipDecompressor()
	zip.RegisterDecompressor(zstd.ZipMethodWinZip, decomp)
	zip.RegisterDecompressor(zstd.ZipMethodPKWare, decomp)
}

func unzstd(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	d, err := zstd.NewReader(in)
	if err != nil {
		return 0, err
	}
	defer d.Close()

	n, err := io.Copy(out, d)
	return n, err
}

func czstd(out io.WriteCloser, in io.ReadCloser, conf any) (int64, error) {
	encoderLvl := zstd.EncoderLevelFromZstd(int(conf.(uint64)))
	log.Debugf("zstd compression level: %d (-> %v)\n", conf.(uint64), encoderLvl)
	e, err := zstd.NewWriter(out, zstd.WithEncoderLevel(encoderLvl))
	if err != nil {
		return 0, err
	}

	n, err := io.Copy(e, in)
	e.Close()
	return n, err
}
