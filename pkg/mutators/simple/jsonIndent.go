package mutators

import (
	"bytes"
	"encoding/json"
	"io"
)

func init() {
	simplestRegister("j", jsonIndent, withDescription("JSON Re-indent"), withHintLexer("JSON"))
}

func jsonIndent(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	j, err := io.ReadAll(r) // NOT streamable
	if err != nil {
		return 0, err
	}
	var b bytes.Buffer
	if len(j) != 0 {
		err = json.Indent(&b, j, "", "  ")
		if err != nil {
			return 0, err
		}
		b.WriteString("\n")
		_, err = io.Copy(w, bytes.NewReader(b.Bytes()))
		if err != nil {
			return 0, err
		}
	}

	return int64(len(j)), nil
}
