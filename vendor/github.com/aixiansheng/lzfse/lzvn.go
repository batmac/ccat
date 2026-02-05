package lzfse

import (
	"bytes"
	"encoding/binary"
	"io"
)

type lzvnDecoder struct {
	r      io.Reader
	w      *cachedWriter
	buffer *bytes.Buffer

	header struct {
		N_raw_bytes     uint32
		N_payload_bytes uint32
	}
}

// at the beginning of each lzvn block
func newLzvnDecoder(r io.Reader, w *cachedWriter) (*lzvnDecoder, error) {
	decoder := &lzvnDecoder{
		r: r,
		w: w,
	}

	if err := binary.Read(r, binary.LittleEndian, &decoder.header); err != nil {
		return nil, err
	}

	b := make([]byte, 0, decoder.header.N_payload_bytes)
	decoder.buffer = bytes.NewBuffer(b)

	return decoder, nil
}

type lzvnOpCode byte

const (
	nop lzvnOpCode = iota
	end_of_stream
	undefined

	small_literal
	large_literal

	small_match
	large_match

	small_distance
	large_distance
	medium_distance
	previous_distance
)

var opcode_table = [256]lzvnOpCode{
	small_distance, small_distance, small_distance, small_distance, small_distance, small_distance, end_of_stream, large_distance,
	small_distance, small_distance, small_distance, small_distance, small_distance, small_distance, nop, large_distance,
	small_distance, small_distance, small_distance, small_distance, small_distance, small_distance, nop, large_distance,
	small_distance, small_distance, small_distance, small_distance, small_distance, small_distance, undefined, large_distance,
	small_distance, small_distance, small_distance, small_distance, small_distance, small_distance, undefined, large_distance,
	small_distance, small_distance, small_distance, small_distance, small_distance, small_distance, undefined, large_distance,
	small_distance, small_distance, small_distance, small_distance, small_distance, small_distance, undefined, large_distance,
	small_distance, small_distance, small_distance, small_distance, small_distance, small_distance, undefined, large_distance,
	small_distance, small_distance, small_distance, small_distance, small_distance, small_distance, previous_distance, large_distance,
	small_distance, small_distance, small_distance, small_distance, small_distance, small_distance, previous_distance, large_distance,
	small_distance, small_distance, small_distance, small_distance, small_distance, small_distance, previous_distance, large_distance,
	small_distance, small_distance, small_distance, small_distance, small_distance, small_distance, previous_distance, large_distance,
	small_distance, small_distance, small_distance, small_distance, small_distance, small_distance, previous_distance, large_distance,
	small_distance, small_distance, small_distance, small_distance, small_distance, small_distance, previous_distance, large_distance,
	undefined, undefined, undefined, undefined, undefined, undefined, undefined, undefined,
	undefined, undefined, undefined, undefined, undefined, undefined, undefined, undefined,
	small_distance, small_distance, small_distance, small_distance, small_distance, small_distance, previous_distance, large_distance,
	small_distance, small_distance, small_distance, small_distance, small_distance, small_distance, previous_distance, large_distance,
	small_distance, small_distance, small_distance, small_distance, small_distance, small_distance, previous_distance, large_distance,
	small_distance, small_distance, small_distance, small_distance, small_distance, small_distance, previous_distance, large_distance,
	medium_distance, medium_distance, medium_distance, medium_distance, medium_distance, medium_distance, medium_distance, medium_distance,
	medium_distance, medium_distance, medium_distance, medium_distance, medium_distance, medium_distance, medium_distance, medium_distance,
	medium_distance, medium_distance, medium_distance, medium_distance, medium_distance, medium_distance, medium_distance, medium_distance,
	medium_distance, medium_distance, medium_distance, medium_distance, medium_distance, medium_distance, medium_distance, medium_distance,
	small_distance, small_distance, small_distance, small_distance, small_distance, small_distance, previous_distance, large_distance,
	small_distance, small_distance, small_distance, small_distance, small_distance, small_distance, previous_distance, large_distance,
	undefined, undefined, undefined, undefined, undefined, undefined, undefined, undefined,
	undefined, undefined, undefined, undefined, undefined, undefined, undefined, undefined,
	large_literal, small_literal, small_literal, small_literal, small_literal, small_literal, small_literal, small_literal,
	small_literal, small_literal, small_literal, small_literal, small_literal, small_literal, small_literal, small_literal,
	large_match, small_match, small_match, small_match, small_match, small_match, small_match, small_match,
	small_match, small_match, small_match, small_match, small_match, small_match, small_match, small_match,
}

type lzvnOp struct {
	code lzvnOpCode
	l    byte
	m    byte
	d    uint16
}

func newLzvnOp(r io.Reader, op *lzvnOp) (*lzvnOp, error) {
	if nil == op {
		op = &lzvnOp{}
	}

	var first_byte byte
	if err := binary.Read(r, binary.LittleEndian, &first_byte); err != nil {
		return nil, err
	}

	op.code = opcode_table[first_byte]
	switch op.code {
	case small_distance:
		var second_byte byte
		if err := binary.Read(r, binary.LittleEndian, &second_byte); err != nil {
			return nil, err
		}

		op.l = extractByte(first_byte, 6, 2)
		op.m = extractByte(first_byte, 3, 3) + 3
		op.d = (extractUint16(uint16(first_byte), 0, 3) << 8) | uint16(second_byte)

	case medium_distance:
		var bytes2and3 uint16
		if err := binary.Read(r, binary.LittleEndian, &bytes2and3); err != nil {
			return nil, err
		}

		op.l = extractByte(first_byte, 3, 2)
		op.m = ((extractByte(first_byte, 0, 3) << 2) | byte(extractUint16(bytes2and3, 0, 2)) + 3)
		op.d = extractUint16(bytes2and3, 2, 14)

	case large_distance:
		var bytes2and3 uint16
		if err := binary.Read(r, binary.LittleEndian, &bytes2and3); err != nil {
			return nil, err
		}

		op.l = extractByte(first_byte, 6, 2)
		op.m = extractByte(first_byte, 3, 3) + 3
		op.d = bytes2and3

	case previous_distance:
		op.l = extractByte(first_byte, 6, 2)
		op.m = extractByte(first_byte, 3, 3) + 3

	case small_match:
		op.m = extractByte(first_byte, 0, 4)

	case large_match:
		var second_byte byte
		if err := binary.Read(r, binary.LittleEndian, &second_byte); err != nil {
			return nil, err
		}
		op.m = second_byte + 16

	case small_literal:
		op.l = extractByte(first_byte, 0, 4)

	case large_literal:
		var second_byte byte
		if err := binary.Read(r, binary.LittleEndian, &second_byte); err != nil {
			return nil, err
		}
		op.l = second_byte + 16

	case end_of_stream:
		var b [7]byte
		if n, err := r.Read(b[:]); n != 7 {
			return nil, err
		}
	}

	return op, nil
}

func (dec *lzvnDecoder) Decode() error {
	var op *lzvnOp
	var opErr error

loop:
	for {
		// Get the operation
		if op, opErr = newLzvnOp(dec.r, op); opErr != nil {
			return opErr
		}

		switch op.code {
		// Distance operations
		case small_distance:
			fallthrough
		case medium_distance:
			fallthrough
		case large_distance:
			fallthrough
		case previous_distance:
			dec.copyLiteralAndMatch(op.l, op.m, op.d)

		// Match operations
		case small_match:
			fallthrough
		case large_match:
			dec.copyMatch(op.m, op.d)

		// Literal operations
		case small_literal:
			fallthrough
		case large_literal:
			dec.copyLiteral(op.l)

		// Other operations
		case nop:
		case end_of_stream:
			break loop
		case undefined:
			panic("Bad LZVN opcode")
		}
	}

	return nil
}

func decodeLZVNBlock(r combinedReader, cw *cachedWriter) error {
	if decoder, err := newLzvnDecoder(r, cw); err != nil {
		return err
	} else {
		return decoder.Decode()
	}
}

// copyMatch copies M bytes from output - D to output
func (dec *lzvnDecoder) copyMatch(m byte, d uint16) error {
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

// copyLiteral copies L bytes from src to dst
func (dec *lzvnDecoder) copyLiteral(l byte) error {
	if _, err := io.CopyN(dec.w, dec.r, int64(l)); err != nil {
		return err
	}

	return nil
}

// copyLiteralAndMatch copies a literal and then a match from output - d to output
func (dec *lzvnDecoder) copyLiteralAndMatch(l, m byte, d uint16) error {
	if err := dec.copyLiteral(l); err != nil {
		return err
	}

	return dec.copyMatch(m, d)
}

func extract64(container uint64, lsb, width int) uint64 {
	if width == 64 {
		return container
	}
	return (container >> lsb) & ((1 << width) - 1)
}

func extractByte(container byte, lsb, width int) byte {
	return byte(extract64(uint64(container), lsb, width))
}

func extractUint16(container uint16, lsb, width int) uint16 {
	return uint16(extract64(uint64(container), lsb, width))
}
