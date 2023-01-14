package mutators

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

func init() {
	singleRegister("j", jsonIndent, withDescription("JSON Re-indent (X:2 space-based)"), withHintLexer("JSON"),
		withConfigBuilder(stdConfigUint64WithDefault(2)),
	)
}

func jsonIndent(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	indent := int(config.(uint64))
	j, err := io.ReadAll(r) // NOT streamable
	if err != nil {
		return 0, err
	}
	var b bytes.Buffer
	if len(j) != 0 {
		err = json.Indent(&b, j, "", strings.Repeat(" ", indent))
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
