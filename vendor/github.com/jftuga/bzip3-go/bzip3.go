// Package bzip3 is a pure-Go port of the bzip3 compression algorithm
// (https://github.com/iczelia/bzip3), file-format compatible with the C
// implementation. The codec pipeline is RLE, LZP, a Burrows-Wheeler
// transform, and an order-0 context-mixing arithmetic coder, applied per
// independent block.
package bzip3

import (
	"encoding/binary"

	"github.com/jftuga/bzip3-go/internal/bwt"
	"github.com/jftuga/bzip3-go/internal/cm"
	"github.com/jftuga/bzip3-go/internal/crc"
	"github.com/jftuga/bzip3-go/internal/lzp"
	"github.com/jftuga/bzip3-go/internal/rle"
)

const (
	// MinBlockSize is the smallest supported block size (65 KiB).
	MinBlockSize = 65 * 1024
	// MaxBlockSize is the largest supported block size (511 MiB).
	MaxBlockSize = 511 * 1024 * 1024
	// DefaultBlockSize matches the upstream CLI default (16 MiB).
	DefaultBlockSize = 16 * 1024 * 1024
)

// Bound returns the recommended output buffer size for compressing
// inputSize bytes (upstream bz3_bound).
func Bound(inputSize int) int { return inputSize + inputSize/50 + 32 }

// State is a single-block codec, the equivalent of upstream bz3_state.
// It is not safe for concurrent use; create one State per goroutine.
type State struct {
	blockSize int32
	swap      []byte
	lut       []int32
	sa        []int32 // BWT scratch, allocated on first encode
	lf        []int32 // unBWT biPSI scratch, allocated on first decode
	cm        cm.Coder
	unbwt     *bwt.Unbwt // unBWT bucket tables, allocated on first decode
}

// NewState creates a codec for blocks of at most blockSize bytes.
func NewState(blockSize int32) (*State, error) {
	if blockSize < MinBlockSize || blockSize > MaxBlockSize {
		return nil, ErrInit
	}
	return &State{
		blockSize: blockSize,
		swap:      make([]byte, Bound(int(blockSize))),
		lut:       make([]int32, lzp.DictSize),
	}, nil
}

// BlockSize returns the block size the state was created with.
func (s *State) BlockSize() int32 { return s.blockSize }

func (s *State) saBuf() []int32 {
	if s.sa == nil {
		s.sa = make([]int32, s.blockSize)
	}
	return s.sa
}

func (s *State) lfBuf() []int32 {
	if s.lf == nil {
		// The biPSI array needs one entry per position plus the sentinel.
		s.lf = make([]int32, Bound(int(s.blockSize))+1)
	}
	return s.lf
}

func (s *State) unbwtBuf() *bwt.Unbwt {
	if s.unbwt == nil {
		s.unbwt = new(bwt.Unbwt)
	}
	return s.unbwt
}

func getS32(b []byte) int32     { return int32(binary.LittleEndian.Uint32(b)) }
func putS32(b []byte, v int32)  { binary.LittleEndian.PutUint32(b, uint32(v)) }
func getU32(b []byte) uint32    { return binary.LittleEndian.Uint32(b) }
func putU32(b []byte, v uint32) { binary.LittleEndian.PutUint32(b, v) }

// Model byte bits (back to front): bit 1 lzp/no lzp, bit 2 rle/no rle.
const (
	modelLZP = 2
	modelRLE = 4
)

// EncodeBlock compresses buffer[:dataSize] in place and returns the
// compressed size. len(buffer) must be at least Bound(dataSize) and
// dataSize at most the state's block size.
func (s *State) EncodeBlock(buffer []byte, dataSize int32) (int32, error) {
	if dataSize > s.blockSize {
		return -1, ErrDataTooBig
	}
	if int64(len(buffer)) < int64(Bound(int(dataSize))) {
		return -1, ErrDataSizeTooSmall
	}

	b1, b2 := buffer, s.swap

	crc32 := crc.Sum(1, b1[:dataSize])

	// Small blocks are stored raw: they would not benefit from entropy
	// coding.
	if dataSize < 64 {
		copy(b1[8:8+dataSize], b1[:dataSize])
		putU32(b1[0:], crc32)
		putS32(b1[4:], -1)
		return dataSize + 8, nil
	}

	var model byte
	var lzpSize, rleSize int32

	rleSize = int32(rle.Encode(b1[:dataSize], b2))
	if rleSize < dataSize {
		b1, b2 = b2, b1
		dataSize = rleSize
		model |= modelRLE
	}

	lzpSize = lzp.Compress(b1[:dataSize], b2[:dataSize], s.lut)
	if lzpSize > 0 && lzpSize < dataSize {
		b1, b2 = b2, b1
		dataSize = lzpSize
		model |= modelLZP
	}

	bwtIdx := bwt.Encode(b2[:dataSize], b1[:dataSize], s.saBuf())

	overhead := int32(2) // CRC32 + BWT index
	if model&modelLZP != 0 {
		overhead++
	}
	if model&modelRLE != 0 {
		overhead++
	}

	s.cm.Begin()
	dataSize = int32(s.cm.Encode(b1[overhead*4+1:], b2[:dataSize]))

	putU32(b1[0:], crc32)
	putS32(b1[4:], bwtIdx)
	b1[8] = model

	p := int32(0)
	if model&modelLZP != 0 {
		putS32(b1[9:], lzpSize)
		p++
	}
	if model&modelRLE != 0 {
		putS32(b1[9+4*p:], rleSize)
	}

	if &b1[0] != &buffer[0] {
		copy(buffer, b1[:dataSize+overhead*4+1])
	}

	return dataSize + overhead*4 + 1, nil
}

// DecodeBlock decompresses a block of compressedSize bytes held in buffer,
// in place, and returns the decompressed size (origSize on success).
// len(buffer) must be sufficient for every intermediate stage; buffers of
// Bound(origSize) bytes always suffice.
func (s *State) DecodeBlock(buffer []byte, compressedSize, origSize int32) (int32, error) {
	bufferSize := len(buffer)

	if bufferSize < 9 || bufferSize < int(compressedSize) {
		return -1, ErrDataSizeTooSmall
	}

	crc32 := getU32(buffer)
	bwtIdx := getS32(buffer[4:])

	bound := int32(Bound(int(s.blockSize)))

	if compressedSize > bound || compressedSize < 0 {
		return -1, ErrMalformedHeader
	}

	if bwtIdx == -1 {
		if compressedSize-8 > 64 || compressedSize < 8 {
			return -1, ErrMalformedHeader
		}
		if int(compressedSize)-8 > bufferSize {
			return -1, ErrDataSizeTooSmall
		}
		copy(buffer, buffer[8:compressedSize])
		if crc.Sum(1, buffer[:compressedSize-8]) != crc32 {
			return -1, ErrCRC
		}
		return compressedSize - 8, nil
	}

	model := buffer[8]

	// Bug-compatible with upstream: (model&2)*4 and (model&4)*4 evaluate
	// to 8 and 16, over-requiring the header size slightly.
	needed := 9 + int(model&modelLZP)*4 + int(model&modelRLE)*4
	if bufferSize < needed {
		return -1, ErrDataSizeTooSmall
	}

	lzpSize, rleSize := int32(-1), int32(-1)
	p := int32(0)
	if model&modelLZP != 0 {
		lzpSize = getS32(buffer[9+4*p:])
		p++
	}
	if model&modelRLE != 0 {
		rleSize = getS32(buffer[9+4*p:])
		p++
	}
	p += 2

	compressedSize -= p*4 + 1

	if (model&modelLZP != 0 && (lzpSize > bound || lzpSize < 0)) ||
		(model&modelRLE != 0 && (rleSize > bound || rleSize < 0)) {
		return -1, ErrMalformedHeader
	}

	if origSize > bound || origSize < 0 {
		return -1, ErrMalformedHeader
	}

	// Size that undoing BWT+CM should decompress into.
	var sizeBeforeBwt int32
	switch {
	case model&modelLZP != 0:
		sizeBeforeBwt = lzpSize
	case model&modelRLE != 0:
		sizeBeforeBwt = rleSize
	default:
		sizeBeforeBwt = origSize
	}

	// Every intermediate stage must fit in buffer (walking backwards, the
	// required size may be lzpSize, rleSize or origSize).
	if int(lzpSize) > bufferSize || int(rleSize) > bufferSize || int(origSize) > bufferSize {
		return -1, ErrDataSizeTooSmall
	}

	b1, b2 := buffer, s.swap

	s.cm.Begin()
	cmStart := int(p*4 + 1)
	cmEnd := cmStart + int(compressedSize)
	if cmEnd < cmStart {
		// Negative payload size from a malformed header: decode from an
		// empty stream (upstream reads -1 bytes past the end instead).
		cmEnd = cmStart
	}
	s.cm.Decode(b1[cmStart:cmEnd], b2[:sizeBeforeBwt])
	b1, b2 = b2, b1

	if bwtIdx > sizeBeforeBwt {
		return -1, ErrMalformedHeader
	}

	if err := s.unbwtBuf().Decode(b2[:sizeBeforeBwt], b1[:sizeBeforeBwt], s.lfBuf(), bwtIdx); err != nil {
		return -1, ErrBWT
	}
	b1, b2 = b2, b1

	sizeSrc := sizeBeforeBwt

	if model&modelLZP != 0 {
		sizeSrc = lzp.Decompress(b1[:lzpSize], b2[:bound], s.lut)
		if sizeSrc == -1 {
			return -1, ErrCRC
		}
		// Data that passed the header checks may still try to escape
		// buffer via its LZP stream; refuse it.
		if int(sizeSrc) > bufferSize {
			return -1, ErrDataSizeTooSmall
		}
		b1, b2 = b2, b1
	}

	if model&modelRLE != 0 {
		if err := rle.Decode(b1[:sizeSrc], b2[:origSize]); err != nil {
			return -1, ErrCRC
		}
		sizeSrc = origSize
		b1, b2 = b2, b1
	}

	if sizeSrc > s.blockSize || sizeSrc < 0 {
		return -1, ErrMalformedHeader
	}

	if sizeSrc > 0 && &b1[0] != &buffer[0] {
		copy(buffer, b1[:sizeSrc])
	}

	if crc32 != crc.Sum(1, buffer[:sizeSrc]) {
		return -1, ErrCRC
	}

	return sizeSrc, nil
}
