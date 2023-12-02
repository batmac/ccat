package mutators

import (
	"io"

	"github.com/batmac/ccat/pkg/log"

	"braces.dev/errtrace"
	qp "mime/quotedprintable"
)

func init() {
	singleRegister("unqp", unqp, withDescription("decode quoted-printable data"),
		withCategory("convert"),
	)
	singleRegister("qp", cqp, withDescription("encode quoted-printable data"),
		withCategory("convert"),
	)
}

func unqp(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	d := qp.NewReader(in)
	if d == nil {
		log.Fatal("qp decoder failed to init.")
	}
	n, err := io.Copy(out, d) // streamable
	return n, errtrace.Wrap(err)
}

func cqp(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	h := qp.NewWriter(out)
	if h == nil {
		log.Fatal("qp encoder failed to init.")
	}
	defer h.Close()
	return errtrace.Wrap2(io.Copy(h, in)) // streamable
}
