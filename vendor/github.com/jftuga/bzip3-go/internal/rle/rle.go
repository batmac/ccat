// Package rle implements bzip3's run-length encoding stage (mrlec/mrled in
// upstream libbz3.c). Runs are collapsed only for byte values where doing so
// yields a net gain; a 32-byte bitmap header records which values are
// collapsed.
package rle

import "errors"

// ErrCorrupt is returned by Decode when the encoded data does not expand to
// exactly the expected output length.
var ErrCorrupt = errors.New("rle: corrupt data")

// Encode compresses in into out and returns the number of bytes written.
// out must be large enough for the worst case (len(in)+32 plus slack); the
// caller passes a bz3_bound-sized buffer as upstream does.
func Encode(in, out []byte) int {
	var t [256]int32
	var run int32
	pc := -1
	for _, cb := range in {
		c := int(cb)
		if c == pc {
			run++
			if run%255 != 0 {
				t[c]++
			}
		} else {
			t[c]--
			run = 0
		}
		pc = c
	}

	op := 0
	for i := 0; i < 32; i++ {
		c := 0
		for j := 0; j < 8; j++ {
			if t[i*8+j] > 0 {
				c += 1 << j
			}
		}
		out[op] = byte(c)
		op++
	}

	pc = -1
	run = 0
	// Mirrors the do/while in upstream mrlec: iterate over all input bytes
	// plus a final -1 terminator.
	for ip := 0; ip <= len(in); ip++ {
		c := -1
		if ip < len(in) {
			c = int(in[ip])
		}
		if c == pc {
			run++
		} else if run > 0 && t[pc] > 0 {
			out[op] = byte(pc)
			op++
			for ; run > 255; run -= 255 {
				out[op] = 255
				op++
			}
			out[op] = byte(run - 1)
			op++
			run = 1
		} else {
			for run++; run > 1; run-- {
				out[op] = byte(pc)
				op++
			}
		}
		pc = c
	}

	return op
}

// Decode expands in into out, which must be exactly the original length.
// It mirrors upstream mrled, including its tolerance of truncated run
// headers (stale pc reuse), and reports ErrCorrupt when the output length
// does not match.
func Decode(in, out []byte) error {
	maxin := len(in)
	outlen := len(out)
	if maxin < 32 {
		return ErrCorrupt
	}

	var t [256]int32
	ip := 0
	for i := 0; i < 32; i++ {
		c := in[ip]
		ip++
		for j := 0; j < 8; j++ {
			t[i*8+j] = int32((c >> j) & 1)
		}
	}

	op := 0
	pc := int32(-1)
	for op < outlen && ip < maxin {
		c := in[ip]
		ip++
		if t[c] != 0 {
			run := int32(0)
			for ip < maxin {
				pc = int32(in[ip])
				ip++
				if pc != 255 {
					break
				}
				run += 255
			}
			run += pc + 1
			for ; run > 0 && op < outlen; run-- {
				out[op] = c
				op++
			}
		} else {
			out[op] = c
			op++
		}
	}

	if op != outlen {
		return ErrCorrupt
	}
	return nil
}
