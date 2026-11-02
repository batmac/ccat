package bzip3

import "io"

// The file format matches the upstream CLI: a 9-byte header ("BZ3v1" plus
// the block size) followed by blocks, each prefixed with its compressed and
// original sizes. Unlike the frame format, there is no block count; the
// stream ends at EOF.

// Writer compresses data written to it into the bzip3 file format.
// Close must be called to flush the final block.
type Writer struct {
	w          io.Writer
	state      *State
	buf        []byte
	n          int32 // bytes buffered for the current block
	wroteHdr   bool
	err        error
	blockCount int64
}

// NewWriter returns a Writer using the given block size (DefaultBlockSize
// if blockSize is 0).
func NewWriter(w io.Writer, blockSize int32) (*Writer, error) {
	if blockSize == 0 {
		blockSize = DefaultBlockSize
	}
	state, err := NewState(blockSize)
	if err != nil {
		return nil, err
	}
	return &Writer{
		w:     w,
		state: state,
		buf:   make([]byte, Bound(int(blockSize))),
	}, nil
}

func (w *Writer) writeHeader() error {
	var hdr [9]byte
	copy(hdr[:], frameMagic[:])
	putS32(hdr[5:], w.state.blockSize)
	_, err := w.w.Write(hdr[:])
	w.wroteHdr = true
	return err
}

// Write buffers p, compressing and emitting a block whenever a full block
// size has accumulated.
func (w *Writer) Write(p []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	if !w.wroteHdr {
		if err := w.writeHeader(); err != nil {
			w.err = err
			return 0, err
		}
	}
	total := 0
	for len(p) > 0 {
		n := copy(w.buf[w.n:w.state.blockSize], p)
		w.n += int32(n)
		p = p[n:]
		total += n
		if w.n == w.state.blockSize {
			if err := w.flushBlock(); err != nil {
				w.err = err
				return total, err
			}
		}
	}
	return total, nil
}

func (w *Writer) flushBlock() error {
	origSize := w.n
	w.n = 0
	newSize, err := w.state.EncodeBlock(w.buf, origSize)
	if err != nil {
		return err
	}
	var bh [8]byte
	putS32(bh[0:], newSize)
	putS32(bh[4:], origSize)
	if _, err := w.w.Write(bh[:]); err != nil {
		return err
	}
	if _, err := w.w.Write(w.buf[:newSize]); err != nil {
		return err
	}
	w.blockCount++
	return nil
}

// Close flushes any buffered data as a final block. It does not close the
// underlying writer.
func (w *Writer) Close() error {
	if w.err != nil {
		return w.err
	}
	if !w.wroteHdr {
		if err := w.writeHeader(); err != nil {
			w.err = err
			return err
		}
	}
	if w.n > 0 {
		if err := w.flushBlock(); err != nil {
			w.err = err
			return err
		}
	}
	w.err = errWriterClosed
	return nil
}
