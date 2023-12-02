package mutators

import (
	"bytes"
	"encoding/json"
	"io"

	"braces.dev/errtrace"
	"github.com/titanous/json5"
)

func init() {
	singleRegister("j5j", j5j, withDescription("JSON5 -> JSON"),
		withHintLexer("JSON"),
		withCategory("convert"),
	)
}

func j5j(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	j, err := io.ReadAll(r) // NOT streamable
	if err != nil {
		return 0, errtrace.Wrap(err)
	}
	if len(j) == 0 {
		return 0, nil
	}
	var b bytes.Buffer
	var c any

	if err := json5.Unmarshal(j, &c); err != nil {
		return 0, errtrace.Wrap(err)
	}

	result, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return 0, errtrace.Wrap(err)
	}
	b.Write(result)
	b.WriteString("\n")
	_, err = io.Copy(w, bytes.NewReader(b.Bytes()))
	if err != nil {
		return 0, errtrace.Wrap(err)
	}

	return int64(len(j)), nil
}
