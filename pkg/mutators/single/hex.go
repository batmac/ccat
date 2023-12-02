package mutators

import (
	"bufio"
	"encoding/hex"
	"io"

	"braces.dev/errtrace"
	"github.com/batmac/ccat/pkg/log"
)

func init() {
	singleRegister("hexdump", hexDump, withDescription("dump in hex as xxd"), withHintLexer("hexdump"))
	singleRegister("hex", hexRaw, withDescription("dump in lowercase hex"),
		withCategory("convert"))
	singleRegister("unhex", unhex, withDescription("decode hex, ignore all non-hex chars"),
		withCategory("convert"))
}

func hexDump(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	dumper := hex.Dumper(w)
	n, err := io.Copy(dumper, r) // streamable
	log.Debugf("finished\n")
	defer dumper.Close()
	return n, errtrace.Wrap(err)
}

func hexRaw(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	h := hex.NewEncoder(w)
	return errtrace.Wrap2(io.Copy(h, r)) // streamable
}

func unhex(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	rp, wp := io.Pipe()
	decoder := hex.NewDecoder(rp)
	go func() {
		scanner := bufio.NewScanner(r)
		scanner.Split(bufio.ScanRunes)
		for scanner.Scan() {
			chars := scanner.Bytes()
			// log.Debugf("scan %v", chars)
			if len(chars) == 1 {
				if isBase16Char(chars[0]) {
					// log.Debugf("base16 %c", chars[0])
					_, _ = wp.Write(chars)
				}
			}
		}
		if err := scanner.Err(); err != nil {
			log.Println("unhexing input:", err)
		}
		_ = wp.Close()
	}()
	return errtrace.Wrap2(io.Copy(w, decoder)) // streamable
}

func isBase16Char(c byte) bool {
	switch {
	case c >= '0' && c <= '9':
		return true
	case c >= 'A' && c <= 'F':
		return true
	case c >= 'a' && c <= 'f':
		return true
	default:
		return false
	}
}
