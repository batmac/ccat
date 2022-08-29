package mutators

import (
	"io"
	"strings"

	"github.com/batmac/ccat/pkg/log"
	"github.com/gabriel-vasile/mimetype"
)

func init() {
	simpleRegister("mimetype", mt, withDescription("detect mimetype"))
}

func mt(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	// we want to be able to end early, so we limit the read
	mimetype.SetLimit(1024)
	mtype, err := mimetype.DetectReader(r)
	if err != nil {
		return 0, err
	}
	log.Debugf("detected mimetype is %s (%s)\n", mtype.String(), mtype.Extension())

	return io.Copy(w, strings.NewReader(mtype.String()+"\n"))
}
