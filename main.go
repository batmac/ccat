package main

import (
	"ccat/color"
	"flag"
	"regexp"
	"strings"
)

var (
	argTokens       = flag.String("t", "", "comma-separated list of tokens")
	argInsensitive  = flag.Bool("i", false, "case-insensitive")
	argOnlyMatching = flag.Bool("o", false, "don't display lines without at least one token")
	argRaw          = flag.Bool("r", false, "don't treat tokens as regexps")
	argLineNumber   = flag.Bool("n", false, "number the output lines, starting at 1.")
	argLockIn       = flag.Bool("L", false, "exclusively flock each file before reading")
	argLockOut      = flag.Bool("l", false, "exclusively flock stdout")
	argSplitByWords = flag.Bool("w", false, "read word by word instead of line by line (only works with utf8)")
	argExec         = flag.String("X", "", "command to exec on each file before processing")

	tmap   map[string]color.Color
	tokens []string
)

func main() {
	flag.Parse()
	if len(*argTokens) > 0 {
		if *argRaw {
			*argTokens = regexp.QuoteMeta(*argTokens)
		}
		tokens = strings.Split(*argTokens, ",")
	}

	tmap = make(map[string]color.Color)
	var c color.Color = new(color.ColorANSI)
	for _, s := range tokens {
		c = c.Next()
		tmap[s] = c
	}
	//fmt.Printf("%v\n", tmap)

	fileList := flag.Args()
	if 0 == len(fileList) {
		fileList = []string{"-"}
	}

	setupStdout(*argLockOut)
	for _, path := range fileList {
		processFile(path)
	}
}
