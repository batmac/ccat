package main

import (
	"ccat/color"
	"ccat/globalctx"
	"ccat/highlighter"
	"ccat/log"
	"ccat/mutators"
	"ccat/openers"
	"ccat/term"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	argTokens       = flag.String("t", "", "comma-separated list of tokens")
	argInsensitive  = flag.Bool("i", false, "tokens given with -t are case-insensitive")
	argOnlyMatching = flag.Bool("o", false, "don't display lines without at least one token")
	argRaw          = flag.Bool("r", false, "don't treat tokens as regexps")
	argLineNumber   = flag.Bool("n", false, "number the output lines, starting at 1.")
	argLockIn       = flag.Bool("L", false, "exclusively flock each file before reading")
	argLockOut      = flag.Bool("l", false, "exclusively flock stdout")
	argSplitByWords = flag.Bool("w", false, "read word by word instead of line by line (only works with utf8)")
	argExec         = flag.String("X", "", "command to exec on each file before processing it")
	argBG           = flag.Bool("bg", false, "colorize the background instead of the font")
	argDebug        = flag.Bool("d", false, "debug what we are doing")
	argHuman        = flag.Bool("H", false, "try to do what is needed to help (syntax-highlight, autodetect, etc. TODO)")
	argStyle        = flag.String("S", "", "style to use (only used if -H, look in -h for the list)")
	argFormatter    = flag.String("F", "", "formatter to use (only used if -H, look in -h for the list)")
	argLexer        = flag.String("P", "", "lexer to use (only used if -H, look in -h for the list)")
	argMutator      = flag.String("m", "", "mutator to use")

	tmap   map[string]color.Color
	tokens []string
)

func init() {
	flag.Usage = Usage
	flag.Parse()

	if !*argDebug {
		log.SetDebug(io.Discard)
	} else {
		log.SetDebug(os.Stderr)
	}

}
func main() {

	log.Debugln("STARTING ccat")

	/* log.Printf("runtest\n")
	err := mutator.RunTest("dummy", os.Stdout, os.Stdin)
	if err != nil {
		log.Fatalln(err)
	} */
	if term.IsStdoutTerminal() {
		*argHuman = true
	}

	if len(*argTokens) > 0 {
		log.Debugln("initializing tokens...")

		if *argRaw {
			*argTokens = regexp.QuoteMeta(*argTokens)
		}
		tokens = strings.Split(*argTokens, ",")
		log.Debugf("tokens: %v\n", tokens)

		log.Debugln("initializing colors...")
		tmap = make(map[string]color.Color)
		var c color.Color
		if *argBG {
			c = new(color.ColorANSIbg)
		} else {
			c = new(color.ColorANSI)
		}
		for _, s := range tokens {
			c = c.Next()
			tmap[s] = c
		}
	}

	//fmt.Printf("%v\n", tmap)

	log.Debugln("initializing file list...")
	fileList := flag.Args()
	if 0 == len(fileList) {
		fileList = []string{"-"}
	}
	log.Debugf("files: %v\n", fileList)

	setupStdout(*argLockOut)

	globalctx.Set("fileList", fileList)
	log.Debugln("processing...")
	for _, path := range fileList {
		processFile(path)
	}
}

func Usage() {
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "---")

	fmt.Fprint(os.Stderr, "ccat <files>...\n")
	fmt.Fprint(os.Stderr, " - highlighter (-H):\n")
	fmt.Fprint(os.Stderr, highlighter.Help())
	fmt.Fprintf(os.Stderr, " - openers:\n    %v\n", strings.Join(openers.ListOpeners(), "\n    "))
	fmt.Fprintf(os.Stderr, " - mutators:\n    %v\n", strings.Join(mutators.ListMutators(), "\n    "))
}
