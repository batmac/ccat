package mutators

import (
	"encoding/base64"
	"io"
)

func init() {
	singleRegister("unbase64", base64Decode, withDescription("decode base64"), withCategory("convert"), withAliases("ub64", "unb64"))
	singleRegister("base64", base64Encode, withDescription("encode to base64"), withCategory("convert"), withAliases("b64"))
}

func base64Decode(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	decoder := base64.NewDecoder(base64.StdEncoding, r)
	return io.Copy(w, decoder) // streamable
}

func base64Encode(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	encoder := base64.NewEncoder(base64.StdEncoding, w)
	written, err := io.Copy(encoder, r) // streamable
	encoder.Close()
	// fmt.Fprintln(w, "")
	return written, err
}
