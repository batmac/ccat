//go:build plugins
// +build plugins

package mutators

import (
	"fmt"
	"io"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// Yaegi is a mutator that executes a yaegi script

func init() {
	singleRegister("yaegi", applyYaegi, withDescription("a yaegi script to apply (path as first argument, symbol as second argument)"),
		withConfigBuilder(stdConfigStrings(2)),
		withCategory("plugin"),
	)
}

func applyYaegi(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	scriptPath := config.([]string)[0]
	symbol := config.([]string)[1]

	i := interp.New(interp.Options{})
	if err := i.Use(stdlib.Symbols); err != nil {
		return 0, err
	}
	if _, err := i.EvalPath(scriptPath); err != nil {
		return 0, err
	}

	sym, err := i.Eval(symbol)
	if err != nil {
		return 0, err
	}
	plugin, ok := sym.Interface().(func(io.WriteCloser, io.ReadCloser, any) (int64, error))
	if !ok {
		return 0, fmt.Errorf("symbol '%s' does not implement the correct signature", symbol)
	}
	return plugin(w, r, config)
}
