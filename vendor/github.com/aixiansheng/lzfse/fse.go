package lzfse

import (
	"encoding/binary"
	"errors"
	"io"
	"fmt"
	"math/bits"
)

type inStream struct {
	r	io.ReadSeeker
	idx         int64
	accum       uint64 // output bits
	accum_nbits int32  // number of valid bits in accum
}

// Initialize the inStream so that its accum holds between 56 and 63 bits.
func newInStream(bits int32, r combinedReader) (*inStream, error) {
	idx, _ := r.Seek(0, io.SeekCurrent)
	fs := &inStream{
		r: io.NewSectionReader(r, 0, idx),
		idx: idx,
	}

	if bits != 0 {
		fs.idx -= 8
		r.Seek(fs.idx, io.SeekStart)
		binary.Read(r, binary.LittleEndian, &fs.accum)
		fs.accum_nbits = bits + 64
	} else {
		fs.idx -= 7
		r.Seek(fs.idx, io.SeekStart)
		accum_bytes := make([]byte, 7)
		r.Read(accum_bytes)
		accum_bytes = append(accum_bytes, 0)

		fs.accum = binary.LittleEndian.Uint64(accum_bytes)
		fs.accum_nbits = bits + 56
	}

	if (fs.accum_nbits < 56 || fs.accum_nbits >= 64) ||
		(fs.accum>>fs.accum_nbits) != 0 {
		return nil, fmt.Errorf("Bad accum/bits %d / %d", fs.accum_nbits, fs.accum>>fs.accum_nbits)
	}

	return fs, nil
}

type lmdDecoderEntry struct {
	total_bits uint8
	value_bits uint8
	delta      int16
	vbase      int32
}

type lmdDecoderTable []lmdDecoderEntry

func newLmdDecoderTable(nstates, nsymbols int, freq []uint16, symbol_vbits []uint8, symbol_vbase []int32) lmdDecoderTable {
	// fse_init_value_decoder_table
	table := make(lmdDecoderTable, 0)
	n_clz := bits.LeadingZeros32(uint32(nstates))
	for i := 0; i < nsymbols; i++ {
		f := int(freq[i])
		if 0 == f {
			continue
		}

		k := bits.LeadingZeros32(uint32(f)) - n_clz
		j0 := ((2 * nstates) >> k) - f

		ei := lmdDecoderEntry{
			value_bits: symbol_vbits[i],
			vbase:      symbol_vbase[i],
		}

		for j := 0; j < f; j++ {
			e := ei

			if j < j0 {
				e.total_bits = uint8(k) + e.value_bits
				e.delta = int16(((f + j) << k) - nstates)
			} else {
				e.total_bits = uint8(k-1) + e.value_bits
				e.delta = int16((j - j0) << (k - 1))
			}

			table = append(table, e)
		}
	}

	return table
}

type fseState uint16

type lmdDecoder struct {
	state fseState
	table lmdDecoderTable
}

func newLmdDecoder(state fseState, table lmdDecoderTable) *lmdDecoder {
	return &lmdDecoder{
		state: state,
		table: table,
	}
}

func (d *lmdDecoder) Decode(in *inStream) int32 {
	// fse_value_decode
	entry := d.table[d.state]
	state_and_value_bits := uint32(in.pull(int32(entry.total_bits)))
	d.state = fseState(uint32(entry.delta) + (state_and_value_bits >> entry.value_bits))
	return int32(uint64(entry.vbase) + fse_mask_lsb64(uint64(state_and_value_bits), uint8(entry.value_bits)))
}

type literalDecoderEntry struct {
	k      int8
	symbol uint8
	delta  int16
}

func (e *literalDecoderEntry) toInt32() int32 {
	return int32(e.k) | (int32(e.symbol) << 8) | (int32(e.delta) << 16)
}

type literalDecoderTable []literalDecoderEntry

func newLiteralDecoderTable(nstates, nsymbols int, freq []uint16) (literalDecoderTable, error) {
	table := make(literalDecoderTable, 1024)
	n_clz := bits.LeadingZeros32(uint32(nstates))
	sum_of_freq := 0
	idx := 0

	for i := 0; i < nsymbols; i++ {
		f := int(freq[i])
		if 0 == f {
			continue
		}

		sum_of_freq += int(f)

		if sum_of_freq > nstates {
			return nil, errors.New("sum_of_freq > nstates")
		}

		k := bits.LeadingZeros32(uint32(f)) - n_clz
		j0 := ((2 * nstates) >> k) - f

		for j := 0; j < f; j++ {
			e := literalDecoderEntry{
				symbol: uint8(i),
			}
			if j < j0 {
				e.k = int8(k)
				e.delta = int16(((f + j) << k) - nstates)
			} else {
				e.k = int8(k - 1)
				e.delta = int16((j - j0) << (k - 1))
			}

			table[idx] = e
			idx++
		}
	}

	return table, nil
}

type literalDecoder struct {
	table literalDecoderTable
	state fseState
}

func newLiteralDecoder(state fseState, table literalDecoderTable) *literalDecoder {
	return &literalDecoder{
		table: table,
		state: state,
	}
}

func (d *literalDecoder) Decode(in *inStream) uint8 {
	e := d.table[d.state]
	eint := e.toInt32()
	d.state = fseState(eint>>16) + fseState(in.pull(eint&0xff))
	return uint8(fse_extract_bits(uint64(eint), 8, 8))
}

var lsb64Mtable = [65]uint64{
	0x0000000000000000, 0x0000000000000001, 0x0000000000000003,
	0x0000000000000007, 0x000000000000000f, 0x000000000000001f,
	0x000000000000003f, 0x000000000000007f, 0x00000000000000ff,
	0x00000000000001ff, 0x00000000000003ff, 0x00000000000007ff,
	0x0000000000000fff, 0x0000000000001fff, 0x0000000000003fff,
	0x0000000000007fff, 0x000000000000ffff, 0x000000000001ffff,
	0x000000000003ffff, 0x000000000007ffff, 0x00000000000fffff,
	0x00000000001fffff, 0x00000000003fffff, 0x00000000007fffff,
	0x0000000000ffffff, 0x0000000001ffffff, 0x0000000003ffffff,
	0x0000000007ffffff, 0x000000000fffffff, 0x000000001fffffff,
	0x000000003fffffff, 0x000000007fffffff, 0x00000000ffffffff,
	0x00000001ffffffff, 0x00000003ffffffff, 0x00000007ffffffff,
	0x0000000fffffffff, 0x0000001fffffffff, 0x0000003fffffffff,
	0x0000007fffffffff, 0x000000ffffffffff, 0x000001ffffffffff,
	0x000003ffffffffff, 0x000007ffffffffff, 0x00000fffffffffff,
	0x00001fffffffffff, 0x00003fffffffffff, 0x00007fffffffffff,
	0x0000ffffffffffff, 0x0001ffffffffffff, 0x0003ffffffffffff,
	0x0007ffffffffffff, 0x000fffffffffffff, 0x001fffffffffffff,
	0x003fffffffffffff, 0x007fffffffffffff, 0x00ffffffffffffff,
	0x01ffffffffffffff, 0x03ffffffffffffff, 0x07ffffffffffffff,
	0x0fffffffffffffff, 0x1fffffffffffffff, 0x3fffffffffffffff,
	0x7fffffffffffffff, 0xffffffffffffffff,
}

func fse_mask_lsb64(x uint64, nbits uint8) uint64 {
	return x & lsb64Mtable[nbits]
}

func fse_extract_bits(x uint64, start, nbits int32) uint64 {
	return fse_mask_lsb64(x>>start, uint8(nbits))
}

// pull consumes the specified number of bits, then returns the value of accum shifted by the remainder.
// accum is restricted to the number of remaining accum bits.
func (fs *inStream) pull(bits int32) uint64 {
	if bits < 0 || bits > fs.accum_nbits {
		panic("bad juju")
	}

	fs.accum_nbits -= bits
	result := fs.accum >> fs.accum_nbits
	fs.accum = fse_mask_lsb64(fs.accum, uint8(fs.accum_nbits))
	return result
}

func (in *inStream) flush() error {
	var nbits int32 = (int32(63) - in.accum_nbits) & int32(-8)
	in.idx -= int64(nbits >> 3)

	in.r.Seek(in.idx, io.SeekStart)
	b := make([]byte, 8)
	in.r.Read(b)
	incoming := binary.LittleEndian.Uint64(b)
	//var incoming uint64
	//binary.Read(in.r, binary.LittleEndian, &incoming)

	in.accum = (in.accum << nbits) | fse_mask_lsb64(incoming, uint8(nbits))
	in.accum_nbits += nbits

	if in.accum_nbits < 56 || in.accum_nbits >= 64 || in.accum>>in.accum_nbits != 0 {
		return errors.New("Bad accum")
	}

	return nil
}

func fse_check_freq(table []uint16, number_of_states int) bool {
	sum_of_freq := 0
	for i := 0; i < len(table); i++ {
		sum_of_freq += int(table[i])
	}

	return sum_of_freq > number_of_states
}
