//go:build plugins
// +build plugins

package mutators

import (
	"fmt"
	"io"
	"os"

	"braces.dev/errtrace"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// Yaegi is a mutator that executes a yaegi script

func init() {
	singleRegister("yaegi", applyYaegi, withDescription("a yaegi script to apply (path as first argument, symbol as second argument)"),
		withConfigBuilder(stdConfigStrings(2, 2)),
		withCategory("plugin"),
	)
}

func applyYaegi(w io.WriteCloser, r io.ReadCloser, configUntyped any) (int64, error) {
	config, ok := configUntyped.([]string)
	if !ok {
		return 0, errtrace.Wrap(fmt.Errorf("config is not a string slice"))
	}
	scriptPath := config[0]
	symbol := config[1]

	i := interp.New(interp.Options{
		Stdin:        r,
		Stdout:       w,
		Stderr:       os.Stderr,
		Args:         os.Args,
		Env:          os.Environ(),
		Unrestricted: false,
	})
	if err := i.Use(stdlib.Symbols); err != nil {
		return 0, errtrace.Wrap(err)
	}

	if _, err := i.EvalPath(scriptPath); err != nil {
		return 0, errtrace.Wrap(err)
	}

	sym, err := i.Eval(symbol)
	if err != nil {
		return 0, errtrace.Wrap(err)
	}
	plugin, ok := sym.Interface().(func(io.WriteCloser, io.ReadCloser, any) (int64, error))
	if !ok {
		return 0, errtrace.Wrap(fmt.Errorf("symbol '%s' does not implement the correct signature", symbol))
	}
	return errtrace.Wrap2(plugin(w, r, config[2:]))
}
