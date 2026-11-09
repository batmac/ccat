package bzip3

import "errors"

// Errors mirror the BZ3_ERR_* codes of upstream libbz3.
var (
	// ErrOutOfBounds corresponds to BZ3_ERR_OUT_OF_BOUNDS.
	ErrOutOfBounds = errors.New("bzip3: data index out of bounds")
	// ErrBWT corresponds to BZ3_ERR_BWT.
	ErrBWT = errors.New("bzip3: Burrows-Wheeler transform failed")
	// ErrCRC corresponds to BZ3_ERR_CRC.
	ErrCRC = errors.New("bzip3: CRC32 check failed")
	// ErrMalformedHeader corresponds to BZ3_ERR_MALFORMED_HEADER.
	ErrMalformedHeader = errors.New("bzip3: malformed header")
	// ErrTruncatedData corresponds to BZ3_ERR_TRUNCATED_DATA.
	ErrTruncatedData = errors.New("bzip3: truncated data")
	// ErrDataTooBig corresponds to BZ3_ERR_DATA_TOO_BIG.
	ErrDataTooBig = errors.New("bzip3: too much data")
	// ErrInit corresponds to BZ3_ERR_INIT.
	ErrInit = errors.New("bzip3: initialization failed")
	// ErrDataSizeTooSmall corresponds to BZ3_ERR_DATA_SIZE_TOO_SMALL.
	ErrDataSizeTooSmall = errors.New("bzip3: buffer too small for decoding")
)
