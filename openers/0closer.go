package openers

import (
	"ccat/log"
	"io"
)

func NewReadCloser(r io.Reader, closure func() error) io.ReadCloser {
	if _, ok := r.(io.WriterTo); ok {
		return newCloserWriterTo{
			Reader:  r,
			closure: closure,
		}
	}
	return newCloser{
		Reader:  r,
		closure: closure,
	}
}

type newCloser struct {
	io.Reader
	closure func() error
}
type checkForClose interface {
	Close() error
}

func (c newCloser) Close() error {
	if _, ok := c.Reader.(checkForClose); ok {
		err := c.Reader.(checkForClose).Close()
		if err != nil {
			log.Println(err)
		}
	}
	return c.closure()
}

type newCloserWriterTo struct {
	io.Reader
	closure func() error
}
type checkForWriterTo interface {
	WriteTo(io.Writer) (int64, error)
}

func (c newCloserWriterTo) Close() error {
	if _, ok := c.Reader.(checkForClose); ok {
		err := c.Reader.(checkForClose).Close()
		if err != nil {
			log.Println(err)
		}
	}
	return c.closure()
}
func (c newCloserWriterTo) WriteTo(w io.Writer) (n int64, err error) {
	return c.Reader.(io.WriterTo).WriteTo(w)
}
