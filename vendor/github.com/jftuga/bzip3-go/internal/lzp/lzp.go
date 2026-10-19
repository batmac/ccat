// Package lzp implements bzip3's LZP (Lempel-Ziv + Prediction) stage, a
// byte-exact port of lzp_encode_block/lzp_decode_block from upstream
// libbz3.c. A 2^18-entry context hash table predicts match positions; only
// matches of at least 40 bytes are collapsed.
package lzp

import "encoding/binary"

const (
	// DictBits is the log2 size of the LZP hash table (LZP_DICTIONARY).
	DictBits = 18
	// DictSize is the number of entries in the LZP hash table.
	DictSize = 1 << DictBits

	minMatch  = 40
	matchByte = 0xf2
)

func hash(ctx uint32) uint32 {
	return (ctx>>15 ^ ctx ^ ctx>>3) & (DictSize - 1)
}

func upcast(b []byte, i int) uint32 {
	return binary.LittleEndian.Uint32(b[i:])
}

// Compress encodes in into out using lut as the hash table (len DictSize,
// cleared here). out must be exactly len(in) bytes: LZP output is only kept
// when it is strictly smaller than its input. Returns the encoded size, or
// -1 when the input is too small or the output would not fit.
func Compress(in, out []byte, lut []int32) int32 {
	if len(in) < minMatch+32 {
		return -1
	}
	clear(lut)
	return encodeBlock(in, out, lut)
}

// Decompress decodes in into out (with len(out) as the output limit),
// using lut as the hash table (len DictSize, cleared here). Returns the
// decoded size, or -1 on truncated input.
func Decompress(in, out []byte, lut []int32) int32 {
	if len(in) < 4 {
		return -1
	}
	clear(lut)
	return decodeBlock(in, out, lut)
}

func encodeBlock(in, out []byte, lut []int32) int32 {
	inEnd := len(in)
	outEOB := len(out) - 8
	heur := 0

	i, o := 0, 0
	for k := 0; k < 4; k++ {
		out[o] = in[i]
		o++
		i++
	}

	ctx := uint32(in[i-1]) | uint32(in[i-2])<<8 | uint32(in[i-3])<<16 | uint32(in[i-4])<<24

	for i < inEnd-minMatch-32 && o < outEOB {
		idx := hash(ctx)
		val := lut[idx]
		lut[idx] = int32(i)
		matched := false
		if val > 0 {
			ref := int(val)
			if upcast(in, i+minMatch-4) == upcast(in, ref+minMatch-4) && upcast(in, i) == upcast(in, ref) {
				ok := true
				if heur > i && upcast(in, heur) != upcast(in, ref+(heur-i)) {
					ok = false
				}
				if ok {
					length := 4
					for ; i+length < inEnd-minMatch-32; length += 4 {
						if upcast(in, i+length) != upcast(in, ref+length) {
							break
						}
					}
					if length < minMatch {
						if heur < i+length {
							heur = i + length
						}
					} else {
						if in[i+length] == in[ref+length] {
							length++
						}
						if in[i+length] == in[ref+length] {
							length++
						}
						if in[i+length] == in[ref+length] {
							length++
						}

						i += length
						ctx = uint32(in[i-1]) | uint32(in[i-2])<<8 | uint32(in[i-3])<<16 | uint32(in[i-4])<<24

						out[o] = matchByte
						o++

						length -= minMatch
						for length >= 254 {
							length -= 254
							out[o] = 254
							o++
							if o >= outEOB {
								break
							}
						}
						out[o] = byte(length)
						o++
						matched = true
					}
				}
			}
			if !matched {
				next := in[i]
				out[o] = next
				o++
				i++
				ctx = ctx<<8 | uint32(next)
				if next == matchByte {
					out[o] = 255
					o++
				}
			}
		} else {
			next := in[i]
			out[o] = next
			o++
			i++
			ctx = ctx<<8 | uint32(next)
		}
	}

	ctx = uint32(in[i-1]) | uint32(in[i-2])<<8 | uint32(in[i-3])<<16 | uint32(in[i-4])<<24

	for i < inEnd && o < outEOB {
		idx := hash(ctx)
		val := lut[idx]
		lut[idx] = int32(i)

		next := in[i]
		out[o] = next
		o++
		i++
		ctx = ctx<<8 | uint32(next)
		if next == matchByte && val > 0 {
			out[o] = 255
			o++
		}
	}

	if o >= outEOB {
		return -1
	}
	return int32(o)
}

func decodeBlock(in, out []byte, lut []int32) int32 {
	inEnd := len(in)
	outEnd := len(out)

	i, o := 0, 0
	for k := 0; k < 4; k++ {
		out[o] = in[i]
		o++
		i++
	}

	ctx := uint32(out[o-1]) | uint32(out[o-2])<<8 | uint32(out[o-3])<<16 | uint32(out[o-4])<<24

	for i < inEnd && o < outEnd {
		idx := hash(ctx)
		val := lut[idx]
		lut[idx] = int32(o)
		if in[i] == matchByte && val > 0 {
			i++
			// 'i' may have been the last index in the case of untrusted bad data.
			if i == inEnd {
				return -1
			}
			if in[i] != 255 {
				length := int32(minMatch)
				for {
					if i == inEnd {
						return -1
					}
					length += int32(in[i])
					if in[i] != 254 {
						i++
						break
					}
					i++
				}

				ref := int(val)
				oe := o + int(length)
				if oe > outEnd {
					oe = outEnd
				}

				// Byte-by-byte: matches may self-overlap.
				for o < oe {
					out[o] = out[ref]
					o++
					ref++
				}

				ctx = uint32(out[o-1]) | uint32(out[o-2])<<8 | uint32(out[o-3])<<16 | uint32(out[o-4])<<24
			} else {
				i++
				out[o] = matchByte
				o++
				ctx = ctx<<8 | uint32(matchByte)
			}
		} else {
			next := in[i]
			i++
			out[o] = next
			o++
			ctx = ctx<<8 | uint32(next)
		}
	}

	return int32(o)
}
