// Package cm implements bzip3's order-0 context-mixing arithmetic coder, a
// byte-exact port of encode_bytes/decode_bytes from upstream libbz3.c
// (based on the coder outlined in Matt Mahoney's DCE). Two counter arrays
// (C0, C1) make the initial bit prediction, refined by an APM/SSE stage (C2)
// keyed on run state.
package cm

// Coder holds the model state. It is reset with Begin before every block,
// matching upstream begin().
type Coder struct {
	c0 [256]uint16
	c1 [256][256]uint16
	c2 [512][17]uint16
}

// Begin resets the model to its initial state.
func (s *Coder) Begin() {
	for i := range s.c0 {
		s.c0[i] = 1 << 15
	}
	for i := range s.c1 {
		for j := range s.c1[i] {
			s.c1[i][j] = 1 << 15
		}
	}
	for i := 0; i < 2; i++ {
		for j := 0; j < 256; j++ {
			for k := 0; k < 17; k++ {
				v := k << 12
				if k == 16 {
					v--
				}
				s.c2[2*j+i][k] = uint16(v)
			}
		}
	}
}

// Encode compresses buf into out and returns the number of bytes written.
// out must be large enough for the worst case; the caller passes a
// bz3_bound-sized region as upstream does.
func (s *Coder) Encode(out, buf []byte) int {
	var high, low uint32 = 0xFFFFFFFF, 0
	c1, c2 := 0, 0
	run := 0
	op := 0

	for _, b := range buf {
		if c1 == c2 {
			run++
		} else {
			run = 0
		}
		f := 0
		if run > 2 {
			f = 1
		}

		ctx := 1
		c := uint32(b)
		c1row := &s.c1[c1]
		c2row := &s.c1[c2]

		for ctx < 256 {
			p0 := int(s.c0[ctx])
			p1 := int(c1row[ctx])
			p2 := int(c2row[ctx])
			p := ((p0+p1)*7 + p2 + p2) >> 4

			j := p >> 12
			apm := &s.c2[2*ctx+f]
			x1 := int(apm[j])
			x2 := int(apm[j+1])
			ssep := x1 + ((x2-x1)*(p&4095))>>12

			if c&128 != 0 {
				high = low + uint32((uint64(high-low)*uint64(ssep*3+p))>>18)

				for low^high < 1<<24 {
					out[op] = byte(low >> 24)
					op++
					low <<= 8
					high = high<<8 + 0xFF
				}

				s.c0[ctx] = uint16(p0 + ((p0 ^ 65535) >> 2))
				c1row[ctx] = uint16(p1 + ((p1 ^ 65535) >> 4))
				apm[j] = uint16(x1 + ((x1 ^ 65535) >> 6))
				apm[j+1] = uint16(x2 + ((x2 ^ 65535) >> 6))
				ctx += ctx + 1
			} else {
				low += uint32((uint64(high-low)*uint64(ssep*3+p))>>18) + 1

				for low^high < 1<<24 {
					out[op] = byte(low >> 24)
					op++
					low <<= 8
					high = high<<8 + 0xFF
				}

				s.c0[ctx] = uint16(p0 - p0>>2)
				c1row[ctx] = uint16(p1 - p1>>4)
				apm[j] = uint16(x1 - x1>>6)
				apm[j+1] = uint16(x2 - x2>>6)
				ctx += ctx
			}

			c = (c << 1) & 0xFF
		}

		c2 = c1
		c1 = ctx & 255
	}

	for k := 0; k < 4; k++ {
		out[op] = byte(low >> 24)
		op++
		low <<= 8
	}
	return op
}

// Decode decompresses len(dst) bytes from src into dst. Reads past the end
// of src yield 0xFF bytes, matching upstream read_in() returning -1.
func (s *Coder) Decode(src, dst []byte) {
	var high, low uint32 = 0xFFFFFFFF, 0
	var code uint32
	c1, c2 := 0, 0
	run := 0

	inPos, inMax := 0, len(src)
	readIn := func() uint32 {
		if inPos < inMax {
			v := uint32(src[inPos])
			inPos++
			return v
		}
		return 0xFFFFFFFF
	}

	for k := 0; k < 4; k++ {
		code = code<<8 + readIn()
	}

	for i := range dst {
		if c1 == c2 {
			run++
		} else {
			run = 0
		}
		f := 0
		if run > 2 {
			f = 1
		}

		ctx := 1
		c1row := &s.c1[c1]
		c2row := &s.c1[c2]

		for ctx < 256 {
			p0 := int(s.c0[ctx])
			p1 := int(c1row[ctx])
			p2 := int(c2row[ctx])
			p := ((p0+p1)*7 + p2 + p2) >> 4

			j := p >> 12
			apm := &s.c2[2*ctx+f]
			x1 := int(apm[j])
			x2 := int(apm[j+1])
			ssep := x1 + ((x2-x1)*(p&4095))>>12

			mid := low + uint32((uint64(high-low)*uint64(ssep*3+p))>>18)
			bit := code <= mid
			if bit {
				high = mid
			} else {
				low = mid + 1
			}
			for low^high < 1<<24 {
				low <<= 8
				high = high<<8 + 255
				code = code<<8 + readIn()
			}

			if bit {
				s.c0[ctx] = uint16(p0 + ((p0 ^ 65535) >> 2))
				c1row[ctx] = uint16(p1 + ((p1 ^ 65535) >> 4))
				apm[j] = uint16(x1 + ((x1 ^ 65535) >> 6))
				apm[j+1] = uint16(x2 + ((x2 ^ 65535) >> 6))
				ctx += ctx + 1
			} else {
				s.c0[ctx] = uint16(p0 - p0>>2)
				c1row[ctx] = uint16(p1 - p1>>4)
				apm[j] = uint16(x1 - x1>>6)
				apm[j+1] = uint16(x2 - x2>>6)
				ctx += ctx
			}
		}

		c2 = c1
		c1 = ctx & 255
		dst[i] = byte(c1)
	}
}
