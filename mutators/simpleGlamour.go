//go:build !nomd
// +build !nomd

package mutators

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/charmbracelet/glamour"
)

func init() {
	// we want the output not to be highlighted again
	simpleRegister("md", glamourize, withDescription("Render Markdown (with glamour)"),
		withExpectingBinary(true))
}

func glamourize(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	g, err := glamour.NewTermRenderer(
		// detect background color and pick either the default dark or light theme
		// glamour.WithAutoStyle(),
		glamour.WithEmoji(),
		glamour.WithEnvironmentConfig(),
	)
	if err != nil {
		return 0, err
	}
	t, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}

	o, err := g.RenderBytes(t)
	if err != nil {
		return 0, err
	}

	return io.Copy(w, bytes.NewReader(o))
}
