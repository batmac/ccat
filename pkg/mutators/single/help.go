package mutators

import (
	"io"
	"strings"

	"braces.dev/errtrace"
	"github.com/batmac/ccat/pkg/mutators"
)

func init() {
	singleRegister("help", printHelp, withDescription("display mutators help"),
		withHintLexer("YAML"),
	)
}

func printHelp(out io.WriteCloser, _ io.ReadCloser, _ any) (int64, error) {
	return errtrace.Wrap2(io.Copy(out, strings.NewReader(mutators.AvailableMutatorsHelp())))
}
