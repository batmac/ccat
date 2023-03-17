package mutators

import (
	"io"
	"strings"
)

func init() {
	// we want the output to be as-is
	singleRegister("filterUTF8", filterUTF8, withDescription("remove non-utf8"),
		withCategory("filter"),
		withExpectingBinary())
}

func filterUTF8(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	u, err := io.ReadAll(r) // NOT streamable
	if err != nil {
		return 0, err
	}

	s := strings.ToValidUTF8(string(u), "")

	return io.Copy(w, strings.NewReader(s))
}
