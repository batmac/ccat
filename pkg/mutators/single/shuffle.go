package mutators

import (
	"fmt"
	"io"
)

// shuffleTargetBlock is the amount of data transposed at once, in bytes. It is
// rounded down to a whole number of elements, and both shuffle and unshuffle
// derive the same value from the element size, so no header is needed to
// round-trip: each side simply reads one block at a time.
const shuffleTargetBlock = 64 * 1024

func init() {
	singleRegister("shuffle", shuffle,
		withDescription("group the X:8-byte elements by byte position, so a following compressor sees runs instead of interleaved bytes"),
		withCategory("filter"),
		withConfigBuilder(stdConfigUint64WithDefault(8)),
	)
	singleRegister("unshuffle", unshuffle,
		withDescription("reverse shuffle (X:8 must match the shuffle element size)"),
		withCategory("filter"),
		withConfigBuilder(stdConfigUint64WithDefault(8)),
	)
}

func shuffle(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	return transposeBlocks(w, r, config, transpose)
}

func unshuffle(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	return transposeBlocks(w, r, config, untranspose)
}

// transposeBlocks streams r to w, applying f to one block at a time.
func transposeBlocks(w io.WriteCloser, r io.ReadCloser, config any, f func(dst, src []byte, n int)) (int64, error) {
	n := cfgInt(config)
	if n < 1 {
		return 0, fmt.Errorf("element size must be at least 1, got %d", n)
	}

	blockSize := shuffleTargetBlock / n * n
	if blockSize == 0 {
		// an element larger than the target block still gets its own block
		blockSize = n
	}
	src := make([]byte, blockSize)
	dst := make([]byte, blockSize)

	var written int64
	for {
		read, err := io.ReadFull(r, src)
		if read > 0 {
			f(dst[:read], src[:read], n)
			nw, werr := w.Write(dst[:read])
			written += int64(nw)
			if werr != nil {
				return written, werr
			}
		}
		switch err {
		case nil:
		case io.EOF, io.ErrUnexpectedEOF:
			return written, nil
		default:
			return written, err
		}
	}
}

// transpose groups src by byte position: every first byte of an element, then
// every second one, and so on. Trailing bytes that do not fill an element are
// copied as-is, and untranspose leaves them alone too, so they round-trip.
func transpose(dst, src []byte, n int) {
	m := len(src) / n
	k := 0
	for j := range n {
		for i := range m {
			dst[k] = src[i*n+j]
			k++
		}
	}
	copy(dst[k:], src[m*n:])
}

// untranspose is the inverse permutation of transpose.
func untranspose(dst, src []byte, n int) {
	m := len(src) / n
	k := 0
	for j := range n {
		for i := range m {
			dst[i*n+j] = src[k]
			k++
		}
	}
	copy(dst[k:], src[m*n:])
}
