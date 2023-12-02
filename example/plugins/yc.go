package plugins

import (
	"braces.dev/errtrace"
	"fmt"
	"io"
	"os"
)

// yc is a simple byte count program meant to be used as a yaegi script.
// (ccat -m yaegi:example/plugins/yc.go:plugins.Y)

func Y(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	count, err := io.Copy(io.Discard, r)
	if err != nil {
		return 0, errtrace.Wrap(err)
	}
	fmt.Fprintln(w, count)
	// dump the entire environment
	/* fmt.Fprintln(w, os.Environ()) */
	// dump the args
	fmt.Println(os.Args)
	// dump the config
	fmt.Println(config)

	return -1, nil
}
