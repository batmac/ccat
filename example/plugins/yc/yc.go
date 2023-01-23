package yc

import (
	"fmt"
	"io"
)

// yc is a simple byte count program meant to be used as a yaegi script.
// (ccat -m yaegi:example/plugins/yc/yc.go:yc.Y)

func Y(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	count, err := io.Copy(io.Discard, r)
	if err != nil {
		return 0, err
	}
	fmt.Fprintln(w, count)

	return -1, nil
}
