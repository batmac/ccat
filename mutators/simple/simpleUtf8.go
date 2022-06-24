package mutators

import (
	"io"
	"io/ioutil"
	"strings"
)

func init() {
	// we want the output to be as-is
	simpleRegister("filterUTF8", filterUTF8, withDescription("remove non-utf8"),
		withCategory("filter"),
		withExpectingBinary(true))
}

func filterUTF8(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	u, err := ioutil.ReadAll(r) // NOT streamable
	if err != nil {
		return 0, err
	}

	s := strings.ToValidUTF8(string(u), "")

	return io.Copy(w, strings.NewReader(s))
}
