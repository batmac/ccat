package mutators

import (
	"bytes"
	"io"

	// "github.com/go-json-experiment/json"
	json "github.com/gowebpki/jcs"
)

func init() {
	simpleRegister("jcs", jcs, withDescription("JSON -> JSON Canonicalization (RFC 8785)"),
		withHintLexer("JSON"),
		withCategory("convert"),
	)
}

func jcs(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	j, err := io.ReadAll(r) // NOT streamable
	if err != nil {
		return 0, err
	}
	if len(j) == 0 {
		return 0, nil
	}
	// var c json.RawValue = j
	// err = c.Canonicalize()
	c, err := json.Transform(j)
	if err != nil {
		return 0, err
	}

	return io.Copy(w, bytes.NewReader(c))
}
