package mutators

import (
	"encoding/hex"
	"io"
)

func init() {
	simpleRegister("hex", "dump in Hex", hexDump)
}

func hexDump(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	dumper := hex.Dumper(w)
	defer dumper.Close()
	return io.Copy(dumper, r)
}
