package mutators

import (
	"fmt"
	"io"
	"strings"

	"howett.net/plist"
	"sigs.k8s.io/yaml"
)

func init() {
	singleNoConfRegister("plist2Y", unplist, withDescription("display an Apple plist as yaml"),
		withHintLexer("YAML"),
		withCategory("convert"),
	)
}

func unplist(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	d, err := io.ReadAll(in) // NOT streamable
	if err != nil {
		return 0, err
	}
	var data any
	_, err = plist.Unmarshal(d, &data)
	if err != nil {
		return 0, err
	}
	d, err = yaml.Marshal(data)
	if err != nil {
		return 0, err
	}
	return io.Copy(out, strings.NewReader(fmt.Sprintf("%v\n", string(d))))
}
