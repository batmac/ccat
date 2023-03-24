package shaped

import (
	"io"
	"time"
)

type Reader struct {
	lastUpdate  time.Time
	r           io.Reader
	bandwidth   int // Bytes per second
	freeVoucher int
	totalRead   int
}

// NewReader creates a new Reader.
// bandwidth is the maximum of bytes to read per second.
func NewReader(r io.Reader, bandwidth int) *Reader {
	return &Reader{
		r:          r,
		bandwidth:  bandwidth,
		lastUpdate: time.Now(),
	}
}

func (lr *Reader) Read(p []byte) (int, error) {
	now := time.Now()
	elapsed := now.Sub(lr.lastUpdate).Seconds()
	lr.lastUpdate = now

	// Calculate the maximum number of bytes that can be read in the elapsed time
	maxBytes := int(elapsed*float64(lr.bandwidth)+0.5) + lr.freeVoucher
	lr.freeVoucher = 0

	if maxBytes == 0 {
		// for safety, should not occur frequently
		maxBytes = 1
		time.Sleep(time.Second)
	}

	// Limit the number of bytes read if the buffer is smaller than the computed maximum
	if maxBytes > len(p) {
		lr.freeVoucher = maxBytes - len(p)
		maxBytes = len(p)
	}

	// Read from the underlying reader
	n, err := lr.r.Read(p[:maxBytes])
	lr.totalRead += n

	return n, err
}
