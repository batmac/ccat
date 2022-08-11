package completion

import (
	_ "embed"
	"log"
	"os"
	"text/template"
)

//go:embed  ccat.tmpl
var tmpl string

type completionData struct {
	Options, Mutators, Formatters, Styles, Lexers string
}

func Print(shell string, opts []string) {
	switch shell {
	case "bash":
		printBash(opts)
	default:
		log.Fatal("completion is currently only available for bash")
	}
}

func printBash(opts []string) {
	data := getCompletionData(opts)

	c, err := template.New("completion").Parse(tmpl)
	if err != nil {
		panic(err)
	}
	err = c.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
