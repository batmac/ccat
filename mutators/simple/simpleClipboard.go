package mutators

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/atotto/clipboard"
	"github.com/batmac/ccat/log"
	"github.com/batmac/ccat/term"
	"github.com/batmac/ccat/utils"
)

func init() {
	simpleRegister("cb", teeClipboard, withDescription("put a copy in the clipboard"))
}

func teeClipboard(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	d, err := ioutil.ReadAll(r) // NOT streamable
	if err != nil {
		return 0, err
	}
	log.Debugf("readall %d bytes\n", len(d))
	if term.IsSsh() || utils.IsRunningInContainer() {
		term.Osc52(d)
	} else {
		cbLocal(d)
	}
	return io.Copy(w, bytes.NewReader(d))
}

func cbLocal(d []byte) {
	log.Debugf("writing to local clipboard\n")
	if err := clipboard.WriteAll(string(d)); err != nil {
		log.Debugln(err)
	}
}
