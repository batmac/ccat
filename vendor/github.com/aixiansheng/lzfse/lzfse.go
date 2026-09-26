package lzfse

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	LZFSE_ENCODE_L_STATES       = 64
	LZFSE_ENCODE_M_STATES       = 64
	LZFSE_ENCODE_D_STATES       = 256
	LZFSE_ENCODE_LITERAL_STATES = 1024

	LZFSE_MATCHES_PER_BLOCK  = 10000
	LZFSE_LITERALS_PER_BLOCK = 4 * LZFSE_MATCHES_PER_BLOCK

	LZFSE_ENCODE_L_SYMBOLS       = 20
	LZFSE_ENCODE_M_SYMBOLS       = 20
	LZFSE_ENCODE_D_SYMBOLS       = 64
	LZFSE_ENCODE_LITERAL_SYMBOLS = 256

	LZFSE_ENCODE_MAX_L_VALUE = 315
	LZFSE_ENCODE_MAX_M_VALUE = 2359
	LZFSE_ENCODE_MAX_D_VALUE = 262139

	LZFSE_N_FREQ = (LZFSE_ENCODE_L_SYMBOLS +
		LZFSE_ENCODE_M_SYMBOLS +
		LZFSE_ENCODE_D_SYMBOLS +
		LZFSE_ENCODE_LITERAL_SYMBOLS)
)

var l_extra_bits = [LZFSE_ENCODE_L_SYMBOLS]uint8{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 3, 5, 8,
}

var l_base_value = [LZFSE_ENCODE_L_SYMBOLS]int32{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 20, 28, 60,
}

var m_extra_bits = [LZFSE_ENCODE_M_SYMBOLS]uint8{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 5, 8, 11,
}

var m_base_value = [LZFSE_ENCODE_M_SYMBOLS]int32{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 24, 56, 312,
}

var d_extra_bits = [LZFSE_ENCODE_D_SYMBOLS]uint8{
	0, 0, 0, 0, 1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3,
	4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7,
	8, 8, 8, 8, 9, 9, 9, 9, 10, 10, 10, 10, 11, 11, 11, 11,
	12, 12, 12, 12, 13, 13, 13, 13, 14, 14, 14, 14, 15, 15, 15, 15,
}

var d_base_value = [LZFSE_ENCODE_D_SYMBOLS]int32{
	0, 1, 2, 3, 4, 6, 8, 10, 12, 16,
	20, 24, 28, 36, 44, 52, 60, 76, 92, 108,
	124, 156, 188, 220, 252, 316, 380, 444, 508, 636,
	764, 892, 1020, 1276, 1532, 1788, 2044, 2556, 3068, 3580,
	4092, 5116, 6140, 7164, 8188, 10236, 12284, 14332, 16380, 20476,
	24572, 28668, 32764, 40956, 49148, 57340, 65532, 81916, 98300, 114684,
	131068, 163836, 196604, 229372,
}

type lzfseDecoder struct {
	v1Header *lzfseV1Header
	w        *cachedWriter

	literals [LZFSE_LITERALS_PER_BLOCK + 64]byte

	fseInStream *inStream
	lDecoder    *lmdDecoder
	mDecoder    *lmdDecoder
	dDecoder    *lmdDecoder
}

type lzfseV1Header struct {
	n_raw_bytes             uint32
	n_payload_bytes         uint32
	n_literals              uint32
	n_matches               uint32
	n_literal_payload_bytes uint32
	n_lmd_payload_bytes     uint32
	literal_bits            uint32
	literal_state           [4]fseState
	lmd_bits                uint32
	l_state                 fseState
	m_state                 fseState
	d_state                 fseState
	l_freq                  [LZFSE_ENCODE_L_SYMBOLS]uint16
	m_freq                  [LZFSE_ENCODE_M_SYMBOLS]uint16
	d_freq                  [LZFSE_ENCODE_D_SYMBOLS]uint16
	literal_freq            [LZFSE_ENCODE_LITERAL_SYMBOLS]uint16
}

type lzfseV2Header struct {
	N_raw_bytes   uint32
	Packed_fields [3]uint64
	Freq          [2 * LZFSE_N_FREQ]byte
}

func (dec *lzfseDecoder) copyMatch(m, d int) error {
	b := make([]byte, m)
	n, err := dec.w.ReadRelativeToEnd(b, int64(d))
	if err == nil && n != len(b) {
		// There weren't enough bytes in the buffer, so we should repeat them until we fill b.
		// (this is what would happen if there was an overlapped copy)
		for i := 0; i < len(b)-n; i++ {
			b[n+i] = b[i]
		}
	}
	_, err = dec.w.Write(b)
	return err
}

func (dec *lzfseDecoder) Decode() error {
	literalIdx := 0
	var D int32 = -1

	for symbols := dec.v1Header.n_matches; symbols > 0; symbols-- {
		if err := dec.fseInStream.flush(); err != nil {
			return err
		}

		L := dec.lDecoder.Decode(dec.fseInStream)
		M := dec.mDecoder.Decode(dec.fseInStream)
		newD := dec.dDecoder.Decode(dec.fseInStream)
		if newD > 0 {
			D = newD
		}

		//fmt.Printf("0x%.16x L=%d M=%d D=%d\n", len(dec.w.buf), L, M, D)

		// Literals...
		b := make([]byte, L)
		m := copy(b, dec.literals[literalIdx:literalIdx+int(L)])
		if n, err := dec.w.Write(b); n != m {
			return err
		}

		literalIdx += int(L)

		// Matches...
		if err := dec.copyMatch(int(M), int(D)); err != nil {
			return err
		}
	}

	return nil
}

// return the int value and the number of bits
func lzfseDecodeV1FreqValue(bits uint32) (int, uint16) {
	lzfse_freq_nbits_table := [32]byte{
		2, 3, 2, 5, 2, 3, 2, 8, 2, 3, 2, 5, 2, 3, 2, 14,
		2, 3, 2, 5, 2, 3, 2, 8, 2, 3, 2, 5, 2, 3, 2, 14,
	}
	lzfse_freq_value_table := [32]byte{
		0, 2, 1, 4, 0, 3, 1, 255, 0, 2, 1, 5, 0, 3, 1, 255,
		0, 2, 1, 6, 0, 3, 1, 255, 0, 2, 1, 7, 0, 3, 1, 255,
	}

	b := bits & 31 // lower 5 bits
	n := uint16(lzfse_freq_nbits_table[b])

	// Special cases for > 5 bits encoding
	if n == 8 {
		return int(n), uint16(8 + ((bits >> 4) & 0xf))
	}

	if n == 14 {
		return int(n), uint16(24 + ((bits >> 4) & 0x3ff))
	}

	// <= 5 bits encoding from table
	return int(n), uint16(lzfse_freq_value_table[b])
}

func (header *lzfseV1Header) Check() (err error) {
	if header.n_literals > LZFSE_LITERALS_PER_BLOCK {
		err = errors.New("header->n_literals > LZFSE_LITERALS_PER_BLOCK")
	} else if header.n_matches > LZFSE_MATCHES_PER_BLOCK {
		err = errors.New("header->n_matches > LZFSE_MATCHES_PER_BLOCK")
	} else if header.literal_state[0] >= LZFSE_ENCODE_LITERAL_STATES {
		err = errors.New("header->literal_state[0] >= LZFSE_ENCODE_LITERAL_STATES")
	} else if header.literal_state[1] >= LZFSE_ENCODE_LITERAL_STATES {
		err = errors.New("header->literal_state[1] >= LZFSE_ENCODE_LITERAL_STATES")
	} else if header.literal_state[2] >= LZFSE_ENCODE_LITERAL_STATES {
		err = errors.New("header->literal_state[2] >= LZFSE_ENCODE_LITERAL_STATES")
	} else if header.literal_state[3] >= LZFSE_ENCODE_LITERAL_STATES {
		err = errors.New("header->literal_state[3] >= LZFSE_ENCODE_LITERAL_STATES")
	} else if header.l_state >= LZFSE_ENCODE_L_STATES {
		err = errors.New("header->l_state >= LZFSE_ENCODE_L_STATES")
	} else if header.m_state >= LZFSE_ENCODE_M_STATES {
		err = errors.New("header->m_state >= LZFSE_ENCODE_M_STATES")
	} else if header.d_state >= LZFSE_ENCODE_D_STATES {
		err = errors.New("header->d_state >= LZFSE_ENCODE_D_STATES")
	} else if fse_check_freq(header.l_freq[:], LZFSE_ENCODE_L_STATES) {
		err = errors.New("fse_check_freq(header.l_freq) failed")
	} else if fse_check_freq(header.m_freq[:], LZFSE_ENCODE_M_STATES) {
		err = errors.New("fse_check_freq(header.m_freq) failed")
	} else if fse_check_freq(header.d_freq[:], LZFSE_ENCODE_D_STATES) {
		err = errors.New("fse_check_freq(header.d_freq) failed")
	} else if fse_check_freq(header.literal_freq[:], LZFSE_ENCODE_LITERAL_STATES) {
		err = errors.New("fsr_check_freq(header->literal_freq) failed")
	}
	return
}

func newLzfseDecoder(r combinedReader, w *cachedWriter, v1 *lzfseV1Header, headerOffset int) (*lzfseDecoder, error) {
	decoder := &lzfseDecoder{
		v1Header: v1,
		w:        w,
	}

	if err := v1.Check(); err != nil {
		return nil, err
	}

	literalDecoderTable, err := newLiteralDecoderTable(LZFSE_ENCODE_LITERAL_STATES, LZFSE_ENCODE_LITERAL_SYMBOLS, v1.literal_freq[:])
	if err != nil {
		return nil, err
	}

	lDecoderTable := newLmdDecoderTable(LZFSE_ENCODE_L_STATES, LZFSE_ENCODE_L_SYMBOLS, v1.l_freq[:], l_extra_bits[:], l_base_value[:])
	mDecoderTable := newLmdDecoderTable(LZFSE_ENCODE_M_STATES, LZFSE_ENCODE_M_SYMBOLS, v1.m_freq[:], m_extra_bits[:], m_base_value[:])
	dDecoderTable := newLmdDecoderTable(LZFSE_ENCODE_D_STATES, LZFSE_ENCODE_D_SYMBOLS, v1.d_freq[:], d_extra_bits[:], d_base_value[:])

	headerOffset += 4
	r.Seek(int64(v1.n_literal_payload_bytes)-int64(headerOffset), io.SeekCurrent)

	in, err := newInStream(int32(v1.literal_bits), r)
	if err != nil {
		return nil, err
	}

	literalDecoder0 := newLiteralDecoder(v1.literal_state[0], literalDecoderTable)
	literalDecoder1 := newLiteralDecoder(v1.literal_state[1], literalDecoderTable)
	literalDecoder2 := newLiteralDecoder(v1.literal_state[2], literalDecoderTable)
	literalDecoder3 := newLiteralDecoder(v1.literal_state[3], literalDecoderTable)

	for i := 0; i < int(v1.n_literals); i += 4 {
		if err := in.flush(); err != nil {
			return nil, err
		}

		decoder.literals[i+0] = literalDecoder0.Decode(in)
		decoder.literals[i+1] = literalDecoder1.Decode(in)
		decoder.literals[i+2] = literalDecoder2.Decode(in)
		decoder.literals[i+3] = literalDecoder3.Decode(in)
	}

	r.Seek(int64(v1.n_lmd_payload_bytes), io.SeekCurrent)

	in2, err := newInStream(int32(v1.lmd_bits), r)
	if err != nil {
		return nil, err
	}

	decoder.lDecoder = newLmdDecoder(v1.l_state, lDecoderTable)
	decoder.mDecoder = newLmdDecoder(v1.m_state, mDecoderTable)
	decoder.dDecoder = newLmdDecoder(v1.d_state, dDecoderTable)

	decoder.fseInStream = in2

	return decoder, nil
}

func v1HeaderFromV2(headerV2 *lzfseV2Header) (*lzfseV1Header, error) {
	v0 := headerV2.Packed_fields[0]
	v1 := headerV2.Packed_fields[1]
	v2 := headerV2.Packed_fields[2]

	headerV1 := &lzfseV1Header{
		n_raw_bytes:             headerV2.N_raw_bytes,
		n_literals:              extract32(v0, 0, 20),
		n_literal_payload_bytes: extract32(v0, 20, 20),
		literal_bits:            extract32(v0, 60, 3) - 7,
		literal_state: [4]fseState{
			fseState(extract32(v1, 0, 10)),
			fseState(extract32(v1, 10, 10)),
			fseState(extract32(v1, 20, 10)),
			fseState(extract32(v1, 30, 10)),
		},
		n_matches:           extract32(v0, 40, 20),
		n_lmd_payload_bytes: extract32(v1, 40, 20),
		lmd_bits:            extract32(v1, 60, 3) - 7,
		l_state:             fseState(extract32(v2, 32, 10)),
		m_state:             fseState(extract32(v2, 42, 10)),
		d_state:             fseState(extract32(v2, 52, 10)),
	}
	headerV1.n_payload_bytes = headerV1.n_literal_payload_bytes - headerV1.n_lmd_payload_bytes

	// Freq tables
	if extract32(v2, 0, 32) == 0 {
		return headerV1, nil
	}

	var accum uint32 = 0
	accum_nbits := 0
	freq_idx := 0
	freq_idx_max := int(extract32(v2, 0, 32)) - 0x20
	freq := make([]uint16, LZFSE_N_FREQ)

	for i := 0; i < LZFSE_N_FREQ; i++ {
		for freq_idx < freq_idx_max && accum_nbits+8 <= 32 {
			accum |= uint32(headerV2.Freq[freq_idx]) << accum_nbits
			accum_nbits += 8
			freq_idx++
		}

		var nbits int
		nbits, freq[i] = lzfseDecodeV1FreqValue(accum)
		if nbits > accum_nbits {
			return nil, errors.New("nbits > accum_nbits")
		}

		accum >>= nbits
		accum_nbits -= nbits
	}

	// This is the most readable way I could do this...
	copy(headerV1.l_freq[:], freq[0:20])       // LZFSE_ENCODE_L_SYMBOLS
	copy(headerV1.m_freq[:], freq[20:40])      // LZFSE_ENCODE_M_SYMBOLS
	copy(headerV1.d_freq[:], freq[40:104])     // LZFSE_ENCODE_D_SYMBOLS
	copy(headerV1.literal_freq[:], freq[104:]) // LZFSE_ENCODE_LITERAL_SYMBOLS

	if accum_nbits >= 8 || freq_idx != freq_idx_max {
		return nil, fmt.Errorf("accum_nbits (%d) >= 8 || freq_idx (%d) != freq_idx_max (%d)",
			accum_nbits, freq_idx, freq_idx_max)
	}

	return headerV1, nil
}

func decodeCompressedV1Block(r combinedReader, w *cachedWriter) error {
	if decoder, err := newLzfseV1Decoder(r, w); err != nil {
		return err
	} else {
		return decoder.Decode()
	}
}

func decodeCompressedV2Block(r combinedReader, w *cachedWriter) error {
	if decoder, err := newLzfseV2Decoder(r, w); err != nil {
		return err
	} else {
		return decoder.Decode()
	}
}

func newLzfseV1Decoder(r combinedReader, w *cachedWriter) (*lzfseDecoder, error) {
	var v1Header lzfseV1Header
	if err := binary.Read(r, binary.LittleEndian, &v1Header); err != nil {
		return nil, err
	} else {
		return newLzfseDecoder(r, w, &v1Header, 0)
	}
}

func newLzfseV2Decoder(r combinedReader, w *cachedWriter) (*lzfseDecoder, error) {
	startLen, _ := r.Seek(0, io.SeekCurrent)
	var v2Header lzfseV2Header
	if err := binary.Read(r, binary.LittleEndian, &v2Header); err != nil {
		return nil, err
	}

	if v1Header, err := v1HeaderFromV2(&v2Header); err != nil {
		return nil, err
	} else {
		endLen, _ := r.Seek(0, io.SeekCurrent)
		totalSeek := endLen - startLen
		headerSize := int64(extract32(v2Header.Packed_fields[2], 0, 32))
		headerOffset := int(totalSeek - headerSize)
		return newLzfseDecoder(r, w, v1Header, headerOffset)
	}
}

func extract32(container uint64, lsb, width int) uint32 {
	if width == 32 {
		return uint32(container >> lsb)
	}
	return uint32((container >> lsb) & ((1 << width) - 1))
}
