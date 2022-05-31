package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
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
	argLockOut      = flag.Bool("l", false, "exclusively flock stdout ")
	argBG           = flag.Bool("bg", false, "colorize the background instead of the font")

	tmap   map[string]Color
	tokens []string
)

type Color interface {
	Sprint(s string) string
	Next() Color
}

func main() {
	flag.Parse()
	if len(*argTokens) > 0 {
		if *argRaw {
			*argTokens = regexp.QuoteMeta(*argTokens)
		}
		tokens = strings.Split(*argTokens, ",")
	}

	tmap = make(map[string]Color)
	var c Color
	if *argBG {
		c = new(ColorANSIbg)
	} else {
		c = new(ColorANSI)
	}
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

func processFile(path string) {

	file := os.Stdin
	var err error
	if path != "-" {
		file, err = fileOpen(path, *argLockIn)
		if err != nil {
			log.Println(err)
			return
		}
		defer fileClose(file, *argLockIn)
	}

	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		var matched bool
		text := scanner.Text()
		for _, token := range tokens {
			var err error
			var regexpPrefix string
			if *argInsensitive {
				regexpPrefix = "(?i)"
			}
			//fmt.Println("text ", text)
			//fmt.Println("token ", token)
			matched, err = regexp.MatchString(regexpPrefix+token, text)
			if err != nil {
				log.Println(err)
			}
			if matched {
				color := tmap[token]
				text = color.Sprint(text)
				break
			}
		}
		if *argLineNumber {
			fmt.Printf("%d ", lineNumber)
			lineNumber++
		}
		if !*argOnlyMatching || matched && *argOnlyMatching {
			fmt.Println(text)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
