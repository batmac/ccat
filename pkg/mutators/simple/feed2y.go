package mutators

import (
	"bytes"
	"io"

	"github.com/mmcdole/gofeed"
	"sigs.k8s.io/yaml"
)

func init() {
	simplestRegister("feed2y", feed2Y, withDescription("rss/atom/json feed -> YAML"),
		withHintLexer("YAML"),
		withCategory("convert"))
}

func feed2Y(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	// convert feed to yaml
	feed, err := gofeed.NewParser().Parse(r)
	if err != nil {
		return 0, err
	}
	y, err := yaml.Marshal(feed)
	if err != nil {
		return 0, err
	}

	n, err := io.Copy(w, bytes.NewReader(y))
	if err != nil {
		return n, err
	}

	return n, nil
}
