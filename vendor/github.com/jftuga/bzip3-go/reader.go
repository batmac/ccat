package bzip3

import (
	"errors"
	"io"
)

var errWriterClosed = errors.New("bzip3: writer already closed")

// Reader decompresses a bzip3 file-format stream.
type Reader struct {
	r       io.Reader
	state   *State
	buf     []byte
	out     []byte // decoded bytes not yet returned
	readHdr bool
	err     error
}

// NewReader returns a Reader for the bzip3 file format. The stream header
// is read lazily on the first Read call.
func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

func (r *Reader) readHeader() error {
	var hdr [9]byte
	if _, err := io.ReadFull(r.r, hdr[:]); err != nil {
		if errors.Is(err, io.ErrUnexpectedEOF) {
			return ErrMalformedHeader
		}
		return err
	}
	if [5]byte(hdr[:5]) != frameMagic {
		return ErrMalformedHeader
	}
	blockSize := getS32(hdr[5:])
	if blockSize < MinBlockSize || blockSize > MaxBlockSize {
		return ErrMalformedHeader
	}
	state, err := NewState(blockSize)
	if err != nil {
		return err
	}
	r.state = state
	r.buf = make([]byte, Bound(int(blockSize)))
	r.readHdr = true
	return nil
}

// nextBlock reads and decodes one block into r.out. Returns io.EOF at a
// clean end of stream.
func (r *Reader) nextBlock() error {
	var bh [8]byte
	if _, err := io.ReadFull(r.r, bh[:]); err != nil {
		if errors.Is(err, io.ErrUnexpectedEOF) {
			return ErrTruncatedData
		}
		return err // io.EOF: clean block boundary
	}
	newSize := getS32(bh[0:])
	origSize := getS32(bh[4:])
	bound := Bound(int(r.state.blockSize))
	if newSize < 0 || int(newSize) > bound || origSize < 0 || int(origSize) > bound {
		return ErrMalformedHeader
	}
	if _, err := io.ReadFull(r.r, r.buf[:newSize]); err != nil {
		if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
			return ErrTruncatedData
		}
		return err
	}
	dSize, err := r.state.DecodeBlock(r.buf, newSize, origSize)
	if err != nil {
		return err
	}
	// See Decompress: refuse blocks whose decoded size disagrees with the
	// advertised original size instead of emitting stale buffer bytes.
	if dSize != origSize {
		return ErrMalformedHeader
	}
	r.out = r.buf[:origSize]
	return nil
}

// Read implements io.Reader.
func (r *Reader) Read(p []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	if !r.readHdr {
		if err := r.readHeader(); err != nil {
			r.err = err
			return 0, err
		}
	}
	for len(r.out) == 0 {
		if err := r.nextBlock(); err != nil {
			r.err = err
			return 0, err
		}
	}
	n := copy(p, r.out)
	r.out = r.out[n:]
	return n, nil
}
