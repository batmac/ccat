package openers

import (
	"io"
)

func NewReadCloser(r io.Reader, closure func() error) io.ReadCloser {
	return &newCloser{
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
		c.Reader.(checkForClose).Close()
	}
	return c.closure()
}
