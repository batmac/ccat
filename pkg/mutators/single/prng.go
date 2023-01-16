package mutators

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"math/rand"
	"time"

	"github.com/batmac/ccat/pkg/log"
)

// pcg32 implements a 32-bit PCG random number generator.
// See http://www.pcg-random.org/ for details.

func init() {
	singleRegister("prng", pcg32, withDescription("generate endless pcg rand (don't use for crypto)"),
		withCategory("random"),
		withExpectingBinary(true),
		withConfigBuilder(stdConfigUint64WithDefault(0)),
	)
}

type PCG32 struct {
	state uint64
	inc   uint64
}

func (pcg *PCG32) Next() uint32 {
	oldstate := pcg.state
	pcg.state = oldstate*6364136223846793005 + (pcg.inc | 1)
	xorshifted := uint32(((oldstate >> 18) ^ oldstate) >> 27)
	rot := uint32(oldstate >> 59)
	return (xorshifted >> rot) | (xorshifted << ((-rot) & 31))
}

func newPCG32(initState uint64, initSeq uint64) *PCG32 {
	pcg := &PCG32{0, (initSeq << 1) | 1}
	pcg.Next()
	pcg.state += initState
	pcg.Next()
	return pcg
}

func pcg32(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	defer r.Close()
	wb := bufio.NewWriterSize(w, 64*1024)
	defer wb.Flush()
	defer w.Close()
	seed := int64(config.(uint64))
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rand.Seed(seed)
	log.Debugf("seed : %d", seed)
	pcg := newPCG32(rand.Uint64(), rand.Uint64()) // #nosec
	n := int64(0)

	for {
		if err := wb.Flush(); err != nil {
			log.Debugf("flush error: %v", err)
			if errors.Is(err, io.ErrClosedPipe) {
				return n, nil
			}
			return 0, err
		}
		availableBytes := wb.Available() / 4 * 4
		b := wb.AvailableBuffer()
		for i := 0; i < availableBytes; i += 4 {
			b = binary.LittleEndian.AppendUint32(b, pcg.Next())
		}
		m, err := wb.Write(b)
		n += int64(m)
		if err != nil {
			return 0, err
		}
	}
}
