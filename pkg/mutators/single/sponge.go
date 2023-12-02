package mutators

import (
	"bytes"
	"io"

	"braces.dev/errtrace"
	"github.com/batmac/ccat/pkg/log"
)

func init() {
	singleRegister("sponge", sponge, withDescription("soak all input before outputting it."))
}

func sponge(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	d, err := io.ReadAll(r) // NOT streamable (that the point :p)
	if err != nil {
		return 0, errtrace.Wrap(err)
	}
	log.Debugf("soaked %d bytes\n", len(d))
	return errtrace.Wrap2(io.Copy(w, bytes.NewReader(d)))
}
