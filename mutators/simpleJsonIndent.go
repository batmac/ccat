package mutators

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

func init() {
	simpleRegister("j", jsonIndent, withDescription("JSON Re-indent"), withHintLexer("JSON"))
}

func jsonIndent(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	j, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}
	var b bytes.Buffer
	err = json.Indent(&b, j, "", "  ")
	if err != nil {
		return 0, err
	}
	b.WriteString("\n")
	io.Copy(w, bytes.NewReader(b.Bytes()))

	return int64(len(j)), nil
}
