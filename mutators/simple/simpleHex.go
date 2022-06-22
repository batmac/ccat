package mutators

import (
	"encoding/hex"
	"io"

	"github.com/batmac/ccat/log"
)

func init() {
	simpleRegister("hexdump", hexDump, withDescription("dump in hex as xxd"), withHintLexer("hexdump"))
	simpleRegister("hex", hexRaw, withDescription("dump in lowercase hex"),
		withCategory("convert"))
}

func hexDump(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	dumper := hex.Dumper(w)
	n, err := io.Copy(dumper, r)
	log.Debugf("finished\n")
	defer dumper.Close()
	return n, err
}

func hexRaw(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	h := hex.NewEncoder(w)
	return io.Copy(h, r)
}
