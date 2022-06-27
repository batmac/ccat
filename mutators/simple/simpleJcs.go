package mutators

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/go-json-experiment/json"
)

func init() {
	simpleRegister("jcs", jcs, withDescription("JSON Canonicalization (RFC 8785)"), withHintLexer("JSON"))
}

func jcs(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	j, err := ioutil.ReadAll(r) // NOT streamable
	if err != nil {
		return 0, err
	}
	if len(j) == 0 {
		return 0, nil
	}
	var c json.RawValue = j
	err = c.Canonicalize()
	if err != nil {
		return 0, err
	}

	return io.Copy(w, bytes.NewReader(c))
}
