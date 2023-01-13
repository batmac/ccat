package mutators

import (
	"bytes"
	"io"

	"github.com/batmac/ccat/pkg/log"
)

func init() {
	singleNoConfRegister("sponge", sponge, withDescription("soak all input before outputting it."))
}

func sponge(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	d, err := io.ReadAll(r) // NOT streamable (that the point :p)
	if err != nil {
		return 0, err
	}
	log.Debugf("soaked %d bytes\n", len(d))
	return io.Copy(w, bytes.NewReader(d))
}
