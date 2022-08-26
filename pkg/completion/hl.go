//go:build !nohl
// +build !nohl

package completion

import (
	"strings"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/batmac/ccat/pkg/mutators"
	"github.com/batmac/ccat/pkg/utils"

	// needed to init the list
	_ "github.com/batmac/ccat/pkg/mutators/simple"
)

func getCompletionData(opts []string) *completionData {
	lexers, styles, formatters := filter(lexers.Names(true), " '"), filter(styles.Names(), " '"), filter(formatters.Names(), " '")
	utils.SortStringsCaseInsensitive(lexers)
	utils.SortStringsCaseInsensitive(styles)
	utils.SortStringsCaseInsensitive(formatters)
	return &completionData{
		Options:    strings.Join(opts, " "),
		Mutators:   strings.Join(mutators.ListAvailableMutators(), " "),
		Formatters: strings.Join(formatters, " "),
		Styles:     strings.Join(styles, " "),
		Lexers:     strings.Join(lexers, " "),
	}
}

func filter(list []string, chars string) []string {
	var result []string
	for _, elem := range list {
		if !strings.ContainsAny(elem, chars) {
			result = append(result, elem)
		}
	}
	return result
}
