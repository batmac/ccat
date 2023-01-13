package mutators

import (
	"io"
	"strings"

	"github.com/batmac/ccat/pkg/log"
	"github.com/gabriel-vasile/mimetype"
)

func init() {
	singlestRegister("mimetype", mt, withDescription("detect mimetype"))
}

func mt(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	mtype, err := mimetype.DetectReader(io.NopCloser(r))
	if err != nil {
		return 0, err
	}
	log.Debugf("detected mimetype is %s (%s)\n", mtype.String(), mtype.Extension())

	// exhaust reader
	_, err = io.Copy(io.Discard, r)
	if err != nil {
		log.Println("mimetype failed to exhaust its reader:", err)
	}
	return io.Copy(w, strings.NewReader(mtype.String()+"\n"))
}
