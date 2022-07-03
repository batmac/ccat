package mutators

import (
	"io"

	"github.com/batmac/ccat/pkg/log"

	qp "mime/quotedprintable"
)

func init() {
	simpleRegister("unqp", unqp, withDescription("decode quoted-printable data"),
		withCategory("convert"),
	)
	simpleRegister("qp", cqp, withDescription("encode quoted-printable data"),
		withCategory("convert"),
	)
}

func unqp(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	d := qp.NewReader(in)
	if d == nil {
		log.Fatal("qp decoder failed to init.")
	}
	n, err := io.Copy(out, d) // streamable
	return n, err
}

func cqp(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	h := qp.NewWriter(out)
	if h == nil {
		log.Fatal("qp encoder failed to init.")
	}
	defer h.Close()
	return io.Copy(h, in) // streamable
}
