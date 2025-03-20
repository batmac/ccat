package mutators

import (
	"io"

	"github.com/minio/minlz"
)

func init() {
	singleRegister("minlz", cminlz, withDescription("compress to minlz data"),
		withCategory("compress"),
		withConfigBuilder(stdConfigUint64WithDefault(minlz.LevelBalanced)),
	)
	singleRegister("unminlz", unminlz, withDescription("decompress minlz data"),
		withCategory("decompress"),
	)
}

func cminlz(out io.WriteCloser, in io.ReadCloser, c any) (int64, error) {
	compressionLevel := int(c.(uint64))
	d := minlz.NewWriter(out, minlz.WriterLevel(compressionLevel))
	n, err := io.Copy(d, in)
	d.Close()
	return n, err
}

func unminlz(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	d := minlz.NewReader(in)
	n, err := io.Copy(out, d)
	return n, err
}
