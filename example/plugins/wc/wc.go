package main

import (
	"bufio"
	"fmt"
	"io"
)

// wc is a simple word count program meant to be used as a go plugin.
// (ccat -m plugin:example/plugins/wc/wc.so:M)

func M(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) { //nolint
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var n int
	for scanner.Scan() {
		n++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	fmt.Fprintln(w, n)
	return -1, nil
}
