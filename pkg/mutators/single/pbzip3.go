package mutators

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"runtime"
	"strconv"

	"github.com/batmac/ccat/pkg/log"
	"github.com/jftuga/bzip3-go"
)

// pbzip3/punbzip3 run bzip3 blocks through an ordered worker pipeline:
// blocks are independent, so the output is byte-identical to bzip3/unbzip3.
// Peak memory is roughly workers * 10 * block size.

const maxBzip3Workers = 1024

type pbzip3Config struct {
	workers   int
	blockSize int32
}

func init() {
	singleRegister("pbzip3", cpbzip3, withDescription("parallel compress to bzip3 data (X:0 is concurrency, 0 is auto, then X:16777216 is block size in bytes)"),
		withCategory("compress"),
		withConfigBuilder(pbzip3ConfigBuilder),
	)
	singleRegister("punbzip3", punbzip3, withDescription("parallel decompress bzip3 data (X:0 is concurrency, 0 is auto)"),
		withCategory("decompress"),
		withConfigBuilder(punbzip3ConfigBuilder),
		withAliases("unpbzip3"),
	)
}

func parseBzip3Workers(arg string) (int, error) {
	w, err := strconv.Atoi(arg)
	if err != nil {
		return 0, err
	}
	if w < 0 || w > maxBzip3Workers {
		return 0, fmt.Errorf("bzip3 concurrency %d out of range [0, %d]", w, maxBzip3Workers)
	}
	return w, nil
}

func pbzip3ConfigBuilder(args []string) (any, error) {
	if len(args) > 2 {
		return nil, ErrWrongNumberOfArgs(0, 2, len(args))
	}
	var c pbzip3Config
	if len(args) > 0 {
		w, err := parseBzip3Workers(args[0])
		if err != nil {
			return nil, err
		}
		c.workers = w
	}
	bs, err := bzip3BlockSizeConfig(args[min(len(args), 1):])
	if err != nil {
		return nil, err
	}
	c.blockSize = bs.(int32)
	return c, nil
}

func punbzip3ConfigBuilder(args []string) (any, error) {
	switch len(args) {
	case 0:
		return 0, nil
	case 1:
		return parseBzip3Workers(args[0])
	default:
		return nil, ErrWrongNumberOfArgs(0, 1, len(args))
	}
}

func bzip3Workers(w int) int {
	if w == 0 {
		return runtime.GOMAXPROCS(0)
	}
	return w
}

// bzip3Job is one block travelling through the pipeline.
type bzip3Job struct {
	buf  []byte
	in   int32 // size of the data in buf before processing
	orig int32 // decompressed size, from the block header (decode only)
	out  int32 // size of the data in buf after processing
	err  error
	done chan struct{}
}

// runBzip3Pipeline reads blocks with next (io.EOF means the end), processes
// them concurrently on per-worker States and hands them to emit in input
// order. At most 2*workers block buffers are allocated.
func runBzip3Pipeline(workers int, blockSize int32,
	next func(buf []byte) (*bzip3Job, error),
	process func(s *bzip3.State, j *bzip3Job),
	emit func(j *bzip3Job) error,
) error {
	quit := make(chan struct{})
	defer close(quit)

	jobs := make(chan *bzip3Job)
	order := make(chan *bzip3Job, workers)
	free := make(chan []byte, 2*workers)
	bound := bzip3.Bound(int(blockSize))

	for range workers {
		go func() {
			var s *bzip3.State // allocated on first use: small inputs only need one
			for j := range jobs {
				if s == nil {
					s, j.err = bzip3.NewState(blockSize)
				}
				if j.err == nil {
					process(s, j)
				}
				close(j.done)
			}
		}()
	}

	go func() {
		defer close(order)
		defer close(jobs)
		allocated := 0
		for {
			var buf []byte
			select {
			case buf = <-free:
			default:
				if allocated < cap(free) {
					buf = make([]byte, bound)
					allocated++
				} else {
					select {
					case buf = <-free:
					case <-quit:
						return
					}
				}
			}

			j, err := next(buf)
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				j = &bzip3Job{err: err, done: make(chan struct{})}
				close(j.done)
				select {
				case order <- j:
				case <-quit:
				}
				return
			}
			j.buf = buf
			j.done = make(chan struct{})

			select {
			case order <- j:
			case <-quit:
				return
			}
			select {
			case jobs <- j:
			case <-quit:
				return
			}
		}
	}()

	for j := range order {
		<-j.done
		if j.err != nil {
			return j.err
		}
		if err := emit(j); err != nil {
			return err
		}
		free <- j.buf
	}
	return nil
}

func cpbzip3(out io.WriteCloser, in io.ReadCloser, conf any) (int64, error) {
	c := conf.(pbzip3Config)
	workers := bzip3Workers(c.workers)
	log.Debugf("pbzip3 workers: %d, block size: %d\n", workers, c.blockSize)

	var hdr [9]byte
	copy(hdr[:], "BZ3v1")
	putBzip3S32(hdr[5:], c.blockSize)
	if _, err := out.Write(hdr[:]); err != nil {
		return 0, err
	}

	eof := false
	next := func(buf []byte) (*bzip3Job, error) {
		if eof {
			return nil, io.EOF
		}
		n, err := io.ReadFull(in, buf[:c.blockSize])
		switch {
		case errors.Is(err, io.EOF), errors.Is(err, io.ErrUnexpectedEOF):
			eof = true
		case err != nil:
			return nil, err
		}
		if n == 0 {
			return nil, io.EOF
		}
		return &bzip3Job{in: int32(n)}, nil // #nosec G115 -- n <= blockSize
	}
	process := func(s *bzip3.State, j *bzip3Job) {
		j.out, j.err = s.EncodeBlock(j.buf, j.in)
	}
	var read int64
	emit := func(j *bzip3Job) error {
		var bh [8]byte
		putBzip3S32(bh[0:], j.out)
		putBzip3S32(bh[4:], j.in)
		if _, err := out.Write(bh[:]); err != nil {
			return err
		}
		_, err := out.Write(j.buf[:j.out])
		read += int64(j.in)
		return err
	}

	err := runBzip3Pipeline(workers, c.blockSize, next, process, emit)
	return read, err
}

func punbzip3(out io.WriteCloser, in io.ReadCloser, conf any) (int64, error) {
	workers := bzip3Workers(conf.(int))

	var hdr [9]byte
	if _, err := io.ReadFull(in, hdr[:]); err != nil {
		if errors.Is(err, io.EOF) {
			return 0, nil // empty input, like unbzip3
		}
		if errors.Is(err, io.ErrUnexpectedEOF) {
			return 0, bzip3.ErrMalformedHeader
		}
		return 0, err
	}
	if string(hdr[:5]) != "BZ3v1" {
		return 0, bzip3.ErrMalformedHeader
	}
	blockSize := getBzip3S32(hdr[5:])
	if blockSize < bzip3.MinBlockSize || blockSize > bzip3.MaxBlockSize {
		return 0, bzip3.ErrMalformedHeader
	}
	log.Debugf("punbzip3 workers: %d, block size: %d\n", workers, blockSize)
	bound := bzip3.Bound(int(blockSize))

	next := func(buf []byte) (*bzip3Job, error) {
		var bh [8]byte
		if _, err := io.ReadFull(in, bh[:]); err != nil {
			if errors.Is(err, io.EOF) {
				return nil, io.EOF // clean block boundary
			}
			if errors.Is(err, io.ErrUnexpectedEOF) {
				return nil, bzip3.ErrTruncatedData
			}
			return nil, err
		}
		newSize, origSize := getBzip3S32(bh[0:]), getBzip3S32(bh[4:])
		if newSize < 0 || int(newSize) > bound || origSize < 0 || int(origSize) > bound {
			return nil, bzip3.ErrMalformedHeader
		}
		if _, err := io.ReadFull(in, buf[:newSize]); err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				return nil, bzip3.ErrTruncatedData
			}
			return nil, err
		}
		return &bzip3Job{in: newSize, orig: origSize}, nil
	}
	process := func(s *bzip3.State, j *bzip3Job) {
		j.out, j.err = s.DecodeBlock(j.buf, j.in, j.orig)
		// like unbzip3: refuse blocks whose decoded size disagrees with the header
		if j.err == nil && j.out != j.orig {
			j.err = bzip3.ErrMalformedHeader
		}
	}
	var written int64
	emit := func(j *bzip3Job) error {
		n, err := out.Write(j.buf[:j.out])
		written += int64(n)
		return err
	}

	err := runBzip3Pipeline(workers, blockSize, next, process, emit)
	return written, err
}

func getBzip3S32(b []byte) int32 {
	return int32(binary.LittleEndian.Uint32(b)) // #nosec G115 -- two's-complement decode
}

func putBzip3S32(b []byte, v int32) {
	binary.LittleEndian.PutUint32(b, uint32(v)) // #nosec G115 -- two's-complement encode
}
