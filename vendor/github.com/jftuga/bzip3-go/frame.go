package bzip3

// The frame format matches upstream bz3_compress/bz3_decompress: a 13-byte
// header ("BZ3v1", block size, block count) followed by blocks, each
// prefixed with its compressed and original sizes.

var frameMagic = [5]byte{'B', 'Z', '3', 'v', '1'}

// Compress compresses src into a self-contained frame using the given
// block size, mirroring upstream bz3_compress (including its block-size
// clamping for small inputs).
func Compress(src []byte, blockSize int32) ([]byte, error) {
	if int64(blockSize) > int64(len(src)) {
		blockSize = int32(Bound(len(src)))
	}
	if blockSize <= MinBlockSize {
		blockSize = MinBlockSize
	}

	state, err := NewState(blockSize)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, Bound(int(blockSize)))

	nBlocks := len(src) / int(blockSize)
	if len(src)%int(blockSize) != 0 {
		nBlocks++
	}

	out := make([]byte, 0, Bound(len(src))+13)
	var hdr [13]byte
	copy(hdr[:], frameMagic[:])
	putS32(hdr[5:], blockSize)
	putS32(hdr[9:], int32(nBlocks))
	out = append(out, hdr[:]...)

	for i, off := 0, 0; i < nBlocks; i++ {
		size := int(blockSize)
		if i == nBlocks-1 {
			size = len(src) % int(blockSize)
		}
		copy(buf[:size], src[off:off+size])
		outSize, err := state.EncodeBlock(buf, int32(size))
		if err != nil {
			return nil, err
		}
		var bh [8]byte
		putS32(bh[0:], outSize)
		putS32(bh[4:], int32(size))
		out = append(out, bh[:]...)
		out = append(out, buf[:outSize]...)
		off += size
	}

	return out, nil
}

// Decompress decompresses a frame produced by Compress (or upstream
// bz3_compress), mirroring upstream bz3_decompress.
func Decompress(src []byte) ([]byte, error) {
	if len(src) < 13 {
		return nil, ErrMalformedHeader
	}
	if [5]byte(src[:5]) != frameMagic {
		return nil, ErrMalformedHeader
	}
	blockSize := getS32(src[5:])
	nBlocks := getS32(src[9:])
	src = src[13:]

	if blockSize < MinBlockSize || blockSize > MaxBlockSize {
		return nil, ErrInit
	}

	state, err := NewState(blockSize)
	if err != nil {
		return nil, err
	}

	bound := Bound(int(blockSize))
	buf := make([]byte, bound)

	// Pre-scan the block headers to size the output exactly.
	total := 0
	rest := src
	for i := int32(0); i < nBlocks; i++ {
		if len(rest) < 8 {
			return nil, ErrMalformedHeader
		}
		size := getS32(rest)
		if size < 0 || int(size) > bound {
			return nil, ErrMalformedHeader
		}
		if len(rest) < int(size)+8 {
			return nil, ErrTruncatedData
		}
		origSize := getS32(rest[4:])
		// Stricter than upstream bz3_decompress, which does not bound
		// orig_size here and can over-read its block buffer.
		if origSize < 0 || int(origSize) > bound {
			return nil, ErrMalformedHeader
		}
		total += int(origSize)
		rest = rest[8+size:]
	}

	out := make([]byte, 0, total)
	for i := int32(0); i < nBlocks; i++ {
		size := getS32(src)
		origSize := getS32(src[4:])
		copy(buf[:size], src[8:8+size])
		dSize, err := state.DecodeBlock(buf, size, origSize)
		if err != nil {
			return nil, err
		}
		// A crafted block can decode successfully to a size other than the
		// advertised orig_size; upstream would copy stale bytes here.
		if dSize != origSize {
			return nil, ErrMalformedHeader
		}
		out = append(out, buf[:origSize]...)
		src = src[8+size:]
	}

	return out, nil
}
