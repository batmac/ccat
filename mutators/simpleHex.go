package mutators

import (
	"encoding/hex"
	"io"
)

func init() {
	simpleRegister("hexdump", hexDump, withDescription("dump in hex as xxd"), withHintLexer("hexdump"))
	simpleRegister("hex", hexRaw, withDescription("dump in lowercase hex"),
		withCategory("convert"))
}

func hexDump(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	dumper := hex.Dumper(w)
	defer dumper.Close()
	return io.Copy(dumper, r)
}

func hexRaw(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	h := hex.NewEncoder(w)
	return io.Copy(h, r)
}
