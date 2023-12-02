package mutators

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"

	"braces.dev/errtrace"
	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/stringutils"
)

func init() {
	singleRegister("pv", pv, withDescription("copy in to out, printing the total and the bandwidth (like pv) each X:1000 milliseconds on stderr"),
		withConfigBuilder(stdConfigUint64WithDefault(1000)))
}

func pv(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	option := config.(uint64)

	done := make(chan struct{})
	defer close(done)
	var totalWritten atomic.Int64

	go func() {
		var oldTotal int64
		fmt.Fprintln(os.Stderr, "")
		defer log.Debugln("goroutine finished")
		// every <option> second, print how many bytes have been transferred or quit if done is closed
		for {
			select {
			case <-done:
				return
			case <-time.After(time.Duration(option) * time.Millisecond):
				prefix := ""
				if log.DebugIsDiscard == 1 {
					prefix = "\x1b[A\x1b[2K" // go on the previous line and erase the line
				}
				newTotal := totalWritten.Load()
				computedDiff := stringutils.HumanSize((newTotal - oldTotal) * 1000 / int64(option))
				oldTotal = newTotal

				fmt.Fprintf(os.Stderr, "%s%s [%s/s]     \n", prefix, stringutils.HumanSize(newTotal), computedDiff)
			}
		}
	}()

	var read, maxRead int
	var rc, wc int64
	var err error

	var bufResizedTimes int
	firstBufSize := os.Getpagesize()
	maxAllowedBufSize := 10 * firstBufSize
	buf := make([]byte, firstBufSize)

	for {
		read, err = r.Read(buf)
		rc++
		// fmt.Fprintf(os.Stderr, "read: %d bytes, read(): %d\n", read, rc)
		if err != nil && read == 0 {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Fprintf(os.Stderr, "read error: %#v\n", err)
			return totalWritten.Load(), errtrace.Wrap(err)
		}
		if read > maxRead {
			maxRead = read
		}
		for {
			m, err := w.Write(buf[:read])
			wc++
			// fmt.Fprintf(os.Stderr, "written: %d bytes, Write(): %d\n", m, wc)

			if err != nil {
				fmt.Fprintf(os.Stderr, "write error: %v\n", err)

				return totalWritten.Load(), errtrace.Wrap(err)
			}
			read -= m
			totalWritten.Add(int64(m))
			if read == 0 {
				break
			}
		}
		if len(buf) < maxAllowedBufSize && maxRead == len(buf) {
			log.Debugf(" resizing b from %d to %d\n", len(buf), len(buf)*2)
			buf = make([]byte, len(buf)*2)
			bufResizedTimes++
		}
	}

	// print counters
	log.Printf(" TOTAL written: %d, read()/write(): %d/%d, buffer min/max size: %d/%d (grown %d times), biggest read was %d\n",
		totalWritten.Load(), rc, wc, firstBufSize, len(buf), bufResizedTimes, maxRead)

	return totalWritten.Load(), nil
}
