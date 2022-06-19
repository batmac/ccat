package mutators

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"howett.net/plist"
	"sigs.k8s.io/yaml"
)

func init() {
	simpleRegister("plist2Y", unplist, withDescription("display an Apple plist as yaml"), withHintLexer("YAML"))
}

func unplist(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	d, err := ioutil.ReadAll(in)
	if err != nil {
		return 0, err
	}
	var data interface{}
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
