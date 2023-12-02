package mutators

import (
	"bytes"
	"io"

	// "github.com/go-json-experiment/json"
	"braces.dev/errtrace"
	json "github.com/gowebpki/jcs"
)

func init() {
	singleRegister("jcs", jcs, withDescription("JSON -> JSON Canonicalization (RFC 8785)"),
		withHintLexer("JSON"),
		withCategory("convert"),
	)
}

func jcs(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	j, err := io.ReadAll(r) // NOT streamable
	if err != nil {
		return 0, errtrace.Wrap(err)
	}
	if len(j) == 0 {
		return 0, nil
	}
	// var c json.RawValue = j
	// err = c.Canonicalize()
	c, err := json.Transform(j)
	if err != nil {
		return 0, errtrace.Wrap(err)
	}

	return errtrace.Wrap2(io.Copy(w, bytes.NewReader(c)))
}
