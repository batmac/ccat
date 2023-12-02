package mutators

import (
	"io"

	"braces.dev/errtrace"
	"github.com/robert-nix/ansihtml"
)

func init() {
	// we want the output to be as-is
	singleRegister("removeANSI", removeANSI, withDescription("remove ANSI codes"),
		withCategory("filter"),
		withExpectingBinary())
}

func removeANSI(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	p := ansihtml.NewParser(r, w)
	if err := p.Parse(nil); err != nil { // streamable
		return 0, errtrace.Wrap(err)
	}

	return 1, nil
}
