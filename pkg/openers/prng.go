//go:build !fileonly
// +build !fileonly

package openers

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/pcg"
	"github.com/batmac/ccat/pkg/utils"
)

var (
	prngOpenerName        = "prng"
	prngOpenerDescription = "generate endless pcg rand (don't use for crypto) (accept a seed as parameter)"
)

type prngOpener struct {
	name, description string
}

func init() {
	register(&prngOpener{
		name:        prngOpenerName,
		description: prngOpenerDescription,
	})
}

func (f prngOpener) Name() string {
	return f.name
}

func (f prngOpener) Description() string {
	return f.description
}

func (f prngOpener) Open(s string, _ bool) (io.ReadCloser, error) {
	s = strings.TrimPrefix(s, "prng://")

	r, w := io.Pipe()
	wb := bufio.NewWriterSize(w, 64*1024)

	seed := time.Now().UnixNano()
	var err error
	if len(s) != 0 {
		seed, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			log.Debugf("error parsing seed: %v", err)
			return nil, err
		}
	}
	rand.Seed(seed)
	log.Debugf("seed : %d", seed)
	pcg := pcg.NewPCG32(rand.Uint64(), rand.Uint64()) // #nosec
	n := int64(0)
	go func() {
		for {
			if err := wb.Flush(); err != nil {
				log.Debugf("flush error: %v", err)
				if errors.Is(err, io.ErrClosedPipe) {
					return
				}
				return
			}
			availableBytes := wb.Available() / 4 * 4
			b := wb.AvailableBuffer()
			for i := 0; i < availableBytes; i += 4 {
				b = binary.LittleEndian.AppendUint32(b, pcg.Next())
			}
			m, err := wb.Write(b)
			n += int64(m)
			if err != nil {
				return
			}
		}
	}()

	return utils.NewReadCloser(r, func() error {
		_ = wb.Flush()
		_ = w.Close()
		return nil
	}), nil
}

func (f prngOpener) Evaluate(s string) float32 {
	if strings.HasPrefix(s, "prng://") {
		return 0.9
	}
	return 0
}
