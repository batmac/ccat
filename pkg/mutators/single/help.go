package mutators

import (
	"io"
	"strings"

	"github.com/batmac/ccat/pkg/mutators"
)

func init() {
	singleNoConfRegister("help", printHelp, withDescription("display mutators help"),
		withHintLexer("YAML"),
	)
}

func printHelp(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	return io.Copy(out, strings.NewReader(mutators.AvailableMutatorsHelp()))
}
