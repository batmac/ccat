//go:build nohl
// +build nohl

package completion

import (
	"strings"

	"github.com/batmac/ccat/pkg/mutators"

	// needed to init the list
	_ "github.com/batmac/ccat/pkg/mutators/simple"
)

func getCompletionData(opts []string) completionData {
	return completionData{
		Options:    strings.Join(opts, " "),
		Mutators:   strings.Join(mutators.ListAvailableMutators(), " "),
		Formatters: "",
		Styles:     "",
		Lexers:     "",
	}
}
