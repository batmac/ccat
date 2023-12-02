package mutators

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/stringutils"
	"github.com/batmac/ccat/pkg/term"
)

// simple mutators to avoid using pipes and chevrons

func init() {
	singleRegister("discard", discard, withDescription("discard X:0 bytes (0 = all)"),
		withConfigBuilder(stdConfigHumanSizeAsInt64WithDefault(0)),
		withAliases("d"),
	)

	singleRegister("wc", wc, withDescription("count bytes (b, default), runes (r), words (w) or lines (l)"),
		withConfigBuilder(stdConfigStringWithDefault("b")))
}

func discard(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	bytes := config.(int64)

	if bytes == 0 {
		log.Debugf("discarding all bytes\n")
		return io.Copy(io.Discard, r)
	}

	log.Debugf("discarding %d bytes\n", bytes)
	_, err := io.CopyN(io.Discard, r, bytes)
	if err != nil && !errors.Is(err, io.EOF) {
		return 0, err
	}

	return io.Copy(w, r)
}

func wc(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	mode := config.(string)

	scanner := bufio.NewScanner(r)

	var splitFn bufio.SplitFunc
	switch mode {
	case "l":
		splitFn = bufio.ScanLines
	case "w":
		splitFn = bufio.ScanWords
	case "r":
		splitFn = bufio.ScanRunes
	case "b":
		splitFn = bufio.ScanBytes
	default:
		return 0, fmt.Errorf("unknown mode '%s'", mode)
	}
	scanner.Split(splitFn)

	var count int64
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	var c string
	if term.IsStdoutTerminal() {
		c = stringutils.HumanSize(count)
	} else {
		c = fmt.Sprintf("%d", count)
	}

	if _, err := fmt.Fprintf(w, "%s\n", c); err != nil {
		return 0, err
	}
	return count, nil
}
