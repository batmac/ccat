//go:build plugins
// +build plugins

package mutators

import (
	"fmt"
	"io"
	"plugin"

	"github.com/batmac/ccat/pkg/log"
)

// https://pkg.go.dev/plugin

func init() {
	singleRegister("plugin", applyPlugin, withDescription("a go plugin to apply (path as first argument, symbol as second argument)"),
		withConfigBuilder(stdConfigStrings(2)),
		withCategory("plugin"),
	)
}

func applyPlugin(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	pluginPath := config.([]string)[0]
	symbol := config.([]string)[1]
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return 0, err
	}
	log.Debugf("plugin %s opened", pluginPath)
	sym, err := p.Lookup(symbol)
	if err != nil {
		return 0, err
	}
	log.Debugf("symbol '%s' found in plugin %s", symbol, pluginPath)
	plugin, ok := sym.(func(io.WriteCloser, io.ReadCloser, any) (int64, error))
	if !ok {
		return 0, fmt.Errorf("symbol '%s' does not implement the correct signature", symbol)
	}
	log.Debugf("%s is a valid symbol", symbol)
	return plugin(w, r, config)
}
