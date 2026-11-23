// Package crc implements the CRC32 variant used by bzip3: the Castagnoli
// polynomial with a caller-supplied initial value (bzip3 uses 1) and no
// final XOR. This matches crc32sum() in upstream libbz3.c.
package crc

import "encoding/binary"

// Sum updates crc with buf and returns the new value. It processes eight
// bytes per step (slicing-by-8, as hash/crc32 does internally), which is
// equivalent to the byte-at-a-time table walk.
func Sum(crc uint32, buf []byte) uint32 {
	for len(buf) >= 8 {
		crc ^= binary.LittleEndian.Uint32(buf)
		hi := binary.LittleEndian.Uint32(buf[4:])
		crc = tables[7][crc&0xFF] ^ tables[6][(crc>>8)&0xFF] ^
			tables[5][(crc>>16)&0xFF] ^ tables[4][crc>>24] ^
			tables[3][hi&0xFF] ^ tables[2][(hi>>8)&0xFF] ^
			tables[1][(hi>>16)&0xFF] ^ tables[0][hi>>24]
		buf = buf[8:]
	}
	for _, b := range buf {
		crc = tables[0][byte(crc)^b] ^ (crc >> 8)
	}
	return crc
}

var tables = makeTables()

func makeTables() (t [8][256]uint32) {
	// Reversed Castagnoli polynomial, as in upstream's hardcoded table.
	const poly = 0x82F63B78
	for i := range t[0] {
		c := uint32(i)
		for k := 0; k < 8; k++ {
			if c&1 != 0 {
				c = poly ^ (c >> 1)
			} else {
				c >>= 1
			}
		}
		t[0][i] = c
	}
	for i := 1; i < 8; i++ {
		for v := range t[i] {
			c := t[i-1][v]
			t[i][v] = t[0][byte(c)] ^ (c >> 8)
		}
	}
	return t
}
