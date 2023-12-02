package mutators

import (
	"io"
	"strings"

	"braces.dev/errtrace"
	"github.com/batmac/ccat/pkg/log"
	"github.com/gabriel-vasile/mimetype"
)

func init() {
	singleRegister("mimetype", mt, withDescription("detect mimetype"))
}

func mt(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	mtype, err := mimetype.DetectReader(io.NopCloser(r))
	if err != nil {
		return 0, errtrace.Wrap(err)
	}
	log.Debugf("detected mimetype is %s (%s)\n", mtype.String(), mtype.Extension())

	// exhaust reader
	_, err = io.Copy(io.Discard, r)
	if err != nil {
		log.Println("mimetype failed to exhaust its reader:", err)
	}
	return errtrace.Wrap2(io.Copy(w, strings.NewReader(mtype.String()+"\n")))
}
