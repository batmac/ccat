package mutators

import (
	"io"

	"github.com/robert-nix/ansihtml"
)

func init() {
	// we want the output to be as-is
	simpleRegister("removeANSI", toHtml, withDescription("remove ANSI codes"), withExpectingBinary(true))
}

func toHtml(w io.WriteCloser, r io.ReadCloser) (int64, error) {

	p := ansihtml.NewParser(r, w)
	p.Parse(nil)

	return 1, nil
}
