package mutators

import (
	"encoding/base64"
	"fmt"
	"io"
)

func init() {
	simpleRegister("unbase64", base64Decode, withDescription("decode base64"))
	simpleRegister("base64", base64Encode, withDescription("encode base64"))
}

func base64Decode(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	decoder := base64.NewDecoder(base64.StdEncoding, r)
	return io.Copy(w, decoder)
}
func base64Encode(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	encoder := base64.NewEncoder(base64.StdEncoding, w)
	written, err := io.Copy(encoder, r)
	encoder.Close()
	fmt.Fprintln(w, "")
	return written, err
}
