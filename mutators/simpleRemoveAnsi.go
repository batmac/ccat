package mutators

import (
	"io"

	"github.com/robert-nix/ansihtml"
)

func init() {
	// we want the output to be as-is
	simpleRegister("removeANSI", removeANSI, withDescription("remove ANSI codes"),
		withCategory("filter"),
		withExpectingBinary(true))
}

func removeANSI(w io.WriteCloser, r io.ReadCloser) (int64, error) {

	p := ansihtml.NewParser(r, w)
	err := p.Parse(nil)
	if err != nil {
		return 0, err
	}

	return 1, nil
}
