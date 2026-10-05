// Package bwt implements the Burrows-Wheeler transform with the exact
// conventions of libsais_bwt/libsais_unbwt as used by bzip3: the output
// starts with the last input byte, the row corresponding to the full string
// is omitted, and the returned primary index is that row's position plus
// one (equivalently, the position of the sentinel row in the BWT of the
// sentinel-terminated string).
package bwt

import (
	"errors"

	"github.com/jftuga/bzip3-go/internal/sais"
)

// ErrInvalidIndex is returned by Decode for a primary index outside [1, n]
// (or != n when n <= 1), matching libsais_unbwt validation.
var ErrInvalidIndex = errors.New("bwt: invalid primary index")

// Encode computes the BWT of in into out (both length n) using sa as
// suffix-array scratch space (len(sa) >= n). It returns the primary index.
func Encode(out, in []byte, sa []int32) int32 {
	n := len(in)
	if n == 0 {
		return 0
	}
	if n == 1 {
		out[0] = in[0]
		return 1
	}

	sa = sa[:n]
	sais.ComputeSA(in, sa)

	// Layout matches libsais_bwt: out[0] is the last input byte, the
	// SA row equal to 0 is skipped, and its position + 1 is returned.
	out[0] = in[n-1]
	idx := int32(-1)
	j := 1
	for i := 0; i < n; i++ {
		s := sa[i]
		if s == 0 {
			idx = int32(i) + 1
			continue
		}
		out[j] = in[s-1]
		j++
	}
	return idx
}

const (
	alphabetSize = 256
	// fastbitsLog matches libsais UNBWT_FASTBITS: the bigram-position
	// space is downsampled to at most 2^17 fastbits entries.
	fastbitsLog = 17
)

// Unbwt holds the fixed-size scratch tables for the bigram inverse BWT
// (512 KiB). It carries no state between calls and can be reused freely,
// but not concurrently.
type Unbwt struct {
	bucket2  [alphabetSize * alphabetSize]int32
	fastbits [1<<fastbitsLog + 1]uint16
}

// Decode inverts the BWT: in holds the transform (length n), idx is the
// primary index returned by Encode, and out receives the original data
// (length n, not aliasing in). p is scratch space with len(p) >= n+1.
//
// This is a port of libsais_unbwt_init_single/libsais_unbwt_decode_1 (the
// single-threaded, single-index path): positions are bucketed by their
// leading bigram so the backward walk emits two bytes per dependent load.
//
// Any idx in the valid range decodes without out-of-bounds access; corrupt
// data simply yields wrong output, which the block-level CRC catches.
func (u *Unbwt) Decode(out, in []byte, p []int32, idx int32) error {
	n := len(in)
	if n <= 1 {
		if int(idx) != n {
			return ErrInvalidIndex
		}
		if n == 1 {
			out[0] = in[0]
		}
		return nil
	}
	if idx <= 0 || int(idx) > n {
		return ErrInvalidIndex
	}
	index := int(idx)

	shift := uint(0)
	for n>>shift > 1<<fastbitsLog {
		shift++
	}

	// Corrupt input can steer the walk onto the two slots biPSI leaves
	// unwritten; zeroing keeps every reachable position within [0, n].
	p = p[:n+1]
	for i := range p {
		p[i] = 0
	}

	var bucket1 [alphabetSize]int32
	for _, c := range in {
		bucket1[c]++
	}

	bucket2 := &u.bucket2
	for i := range bucket2 {
		bucket2[i] = 0
	}

	// Per first-character range, histogram the BWT character of each row
	// (rows are the sentinel-augmented matrix rows 1..n; row index is the
	// sentinel row and is skipped, rows above it read one position lower).
	sum := int32(1)
	for c := 0; c < alphabetSize; c++ {
		prev := sum
		sum += bucket1[c]
		bucket1[c] = prev
		if prev == sum {
			continue
		}
		b2 := bucket2[c<<8 : c<<8+alphabetSize]
		hi := int(sum)
		if index < hi {
			hi = index
		}
		if hi > int(prev) {
			for _, d := range in[prev:hi] {
				b2[d]++
			}
		}
		lo := index + 1
		if int(prev) > lo {
			lo = int(prev)
		}
		if int(sum) > lo {
			for _, d := range in[lo-1 : sum-1] {
				b2[d]++
			}
		}
	}

	// Transpose so bucket2 is keyed by (first char << 8) | second char.
	for c := 0; c < alphabetSize; c++ {
		for d := c + 1; d < alphabetSize; d++ {
			bucket2[d<<8+c], bucket2[c<<8+d] = bucket2[c<<8+d], bucket2[d<<8+c]
		}
	}

	// Prefix-sum bucket2 into bucket starts and build the fastbits index.
	// The extra slot at lastc reserves the sentinel row's bigram position.
	lastc := int(in[0])
	fastbits := &u.fastbits
	v := 0
	sum = 1
	w := 0
	for c := 0; c < alphabetSize; c++ {
		if c == lastc {
			sum++
		}
		for d := 0; d < alphabetSize; d++ {
			prev := sum
			sum += bucket2[w]
			bucket2[w] = prev
			if prev != sum {
				for ; v <= int(sum-1)>>shift; v++ {
					fastbits[v] = uint16(w)
				}
			}
			w++
		}
	}

	// biPSI: for each text position, place it in its bigram bucket. After
	// this pass each bucket2 entry holds its bucket's end position.
	for i := 0; i < index; i++ {
		c := in[i]
		pp := bucket1[c]
		bucket1[c] = pp + 1
		if int(pp) != index {
			var d byte
			if int(pp) < index {
				d = in[pp]
			} else {
				d = in[pp-1]
			}
			w := int(d)<<8 + int(c)
			p[bucket2[w]] = int32(i)
			bucket2[w]++
		}
	}
	for i := index + 1; i <= n; i++ {
		c := in[i-1]
		pp := bucket1[c]
		bucket1[c] = pp + 1
		if int(pp) != index {
			var d byte
			if int(pp) < index {
				d = in[pp]
			} else {
				d = in[pp-1]
			}
			w := int(d)<<8 + int(c)
			p[bucket2[w]] = int32(i)
			bucket2[w]++
		}
	}

	// Walk forward from the primary index, two output bytes per step.
	p0 := int32(index)
	for i := 0; i < n>>1; i++ {
		c0 := int(fastbits[p0>>shift])
		for bucket2[c0] <= p0 {
			c0++
		}
		p0 = p[p0]
		out[2*i] = byte(c0 >> 8)
		out[2*i+1] = byte(c0)
	}
	// The last byte of the original data is always the BWT's first byte
	// (for even n this overwrites the walk's final byte with the same
	// value on well-formed input).
	out[n-1] = in[0]
	return nil
}
