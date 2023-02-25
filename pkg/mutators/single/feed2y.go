package mutators

import (
	"bytes"
	"io"

	"github.com/mmcdole/gofeed"
	"sigs.k8s.io/yaml"
)

func init() {
	singleRegister("feed2y", feed2Y, withDescription("rss/atom/json feed -> YAML"),
		withHintLexer("YAML"),
		withCategory("convert"))
}

func feed2Y(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
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

	return n, err
}
