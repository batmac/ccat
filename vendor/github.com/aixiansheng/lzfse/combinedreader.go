package lzfse

import (
	"io"
)

type combinedReader interface {
	io.ReaderAt
	io.Reader
	io.ReadSeeker
}
