package plugins

import (
	"braces.dev/errtrace"
	"io"
	"os"
)

// Example: ccat -m yaegi:example/plugins/tee.go:plugins.Tee,sha256 go.mod

func Tee(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	// open a tee.raw file for writing
	file, err := os.Create("tee.raw")
	if err != nil {
		return 0, errtrace.Wrap(err)
	}
	defer file.Close()

	tee := io.TeeReader(r, file)

	return errtrace.Wrap2(io.Copy(w, tee))
}
