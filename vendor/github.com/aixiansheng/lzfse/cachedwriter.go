package lzfse

import (
	"bytes"
)

type cachedWriter struct {
	buf *bytes.Buffer
}

func newCachedWriter() *cachedWriter {
	return &cachedWriter{
		buf: &bytes.Buffer{},
	}
}

func (w *cachedWriter) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}

func (w *cachedWriter) ReadRelativeToEnd(b []byte, offset int64) (copied int, err error) {
	bb := w.buf.Bytes()
	copied = copy(b, bb[int64(w.buf.Len())-offset:])
	return
}

func (w *cachedWriter) Bytes() []byte {
	return w.buf.Bytes()
}
