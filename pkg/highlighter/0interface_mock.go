//go:build nohl
// +build nohl

package highlighter

import (
	"io"
)

func Go(w io.WriteCloser, r io.ReadCloser, _ Options) error {
	go func() {
		_, _ = io.Copy(w, r)
		w.Close()
	}()
	return nil
}

func Help() string {
	return "not supported (compiled with nohl)\n"
}

func Run(input string, _ *Options) string {
	return input
}
