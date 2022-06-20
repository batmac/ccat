package mutators

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"sigs.k8s.io/yaml"
)

func init() {
	simpleRegister("j2y", J2Y, withDescription("JSON -> YAML"),
		withHintLexer("YAML"),
		withCategory("convert"))
	simpleRegister("y2j", Y2J, withDescription("YAML -> JSON"),
		withHintLexer("JSON"),
		withCategory("convert"))
}

func J2Y(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	j, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}
	y, err := yaml.JSONToYAML(j)
	if err != nil {
		return 0, err
	}

	_, err = io.Copy(w, bytes.NewReader(y))
	if err != nil {
		return 0, err
	}

	return int64(len(j)), nil
}

func Y2J(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	y, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}
	j, err := yaml.YAMLToJSON(y)
	if err != nil {
		return 0, err
	}
	var b bytes.Buffer
	err = json.Indent(&b, j, "", "  ")
	if err != nil {
		return 0, err
	}
	b.WriteString("\n")

	_, err = io.Copy(w, bytes.NewReader(b.Bytes()))
	if err != nil {
		return 0, err
	}
	return int64(len(y)), nil
}
