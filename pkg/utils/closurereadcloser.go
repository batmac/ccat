package utils

import (
	"io"

	"github.com/batmac/ccat/pkg/log"
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
	log.Debugf("newCloser Close()\n")
	if _, ok := c.Reader.(checkForClose); ok {
		log.Debugf("newCloser inner Close()\n")
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

func (c newCloserWriterTo) Close() error {
	log.Debugf("newCloserWriterTo Close()\n")
	if _, ok := c.Reader.(checkForClose); ok {
		log.Debugf("newCloserWriterTo inner Close()\n")
		err := c.Reader.(checkForClose).Close()
		if err != nil {
			log.Println(err)
		}
	}
	return c.closure()
}

func (c newCloserWriterTo) WriteTo(w io.Writer) (int64, error) {
	return c.Reader.(io.WriterTo).WriteTo(w)
}
