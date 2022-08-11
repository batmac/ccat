package completion

import (
	_ "embed"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/batmac/ccat/pkg/mutators"
	"github.com/batmac/ccat/pkg/utils"

	// needed to init the list
	_ "github.com/batmac/ccat/pkg/mutators/simple"
)

//go:embed  ccat.tmpl
var tmpl string

type Completion struct {
	Options, Mutators, Formatters, Styles, Lexers string
}

func Print(shell string, opts []string) {
	if shell != "bash" {
		log.Fatal("completion is currently only available for bash")
	}
	lexers, styles, formatters := filter(lexers.Names(true), " '"), filter(styles.Names(), " '"), filter(formatters.Names(), " '")
	utils.SortStringsCaseInsensitive(lexers)
	utils.SortStringsCaseInsensitive(styles)
	utils.SortStringsCaseInsensitive(formatters)

	data := Completion{
		Options:    strings.Join(opts, " "),
		Mutators:   strings.Join(mutators.ListAvailableMutators(), " "),
		Formatters: strings.Join(formatters, " "),
		Styles:     strings.Join(styles, " "),
		Lexers:     strings.Join(lexers, " "),
	}

	c, err := template.New("completion").Parse(tmpl)
	if err != nil {
		panic(err)
	}
	err = c.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
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
