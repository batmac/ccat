package mutators

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/atotto/clipboard"
	"github.com/batmac/ccat/log"
)

func init() {
	simpleRegister("cb", teeClipboard, withDescription("put a copy in the clipboard"))
}

func teeClipboard(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	d, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}
	log.Debugf("readall %d bytes\n", len(d))
	if err := clipboard.WriteAll(string(d)); err != nil {
		log.Debugln(err)
	}

	return io.Copy(w, bytes.NewReader(d))
}
