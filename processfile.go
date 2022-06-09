package main

import (
	"bufio"
	"ccat/highlighter"
	"ccat/log"
	"ccat/mutators"
	"ccat/openers"
	"ccat/pipedcmd"
	"ccat/scanners"
	"ccat/term"
	"fmt"
	"io"

	//_ "net/http/pprof"
	"os"
	"regexp"
)

func processFile(path string) {
	log.Debugf("processing %s...\n", path)

	from, err := openers.Open(path, *argLockIn)
	if err != nil {
		log.Printf("opening %s: %v", path, err)
		return
	}

	if len(*argMutator) > 0 {
		choice := *argMutator
		r, w := io.Pipe()
		m, err := mutators.New(choice)
		if err != nil {
			log.Fatal(err)
		}
		if m.Start(w, from) != nil {
			log.Fatal("failed to start the mutator\n")
		}

		from = r
	}
	fromOrig := from

	defer func() {
		// I don't want to determine if already closed, try to close it, it will fail if it is already closed
		_ = fromOrig.Close()
		log.Debugf("closed %s...\n", path)
	}()
	/*************************************/
	if len(*argExec) > 0 {
		log.Debugf("creating pipedcmd %v...\n", *argExec)
		cmd, err := pipedcmd.New(*argExec)
		//log.Debugf("%s", log.Pp(cmd))

		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			log.Debugf("waiting pipedcmd %v...\n", *argExec)

			if err := cmd.Wait(); err != nil {
				log.Println(err)
			}
		}()

		log.Debugf("start pipedcmd %s\n", cmd)

		err = cmd.Start(from)
		if err != nil {
			log.Println(err)
		}
		from = cmd.Stdout.(*os.File)
	}

	if *argHuman {
		if term.IsArt(path) {
			log.Debugf("is art, displaying directly...\n")
			term.PrintArt(from)
			return
		}
		log.Debugln("highlighting...")
		r, w := io.Pipe()
		err := highlighter.Go(w, from, highlighter.Options{
			FileName:      path,
			StyleHint:     *argStyle,
			FormatterHint: *argFormatter,
			LexerHint:     *argLexer,
		})
		if err != nil {
			log.Printf("error while highlighting: %v", err)
		} else {
			from = r
		}
	}
	log.Debugln("initializing Scanner...")

	//go http.ListenAndServe(":8090", nil)
	scanner := bufio.NewScanner(from)

	splitFn := scanners.ScanBytes
	//splitFn := scanners.ScanLines
	if len(tokens) > 0 || *argLineNumber || *argOnlyMatching {
		log.Debugln("splitting on Lines...")
		splitFn = scanners.ScanLines
	}
	if *argSplitByWords {
		log.Debugln("splitting on Words...")
		splitFn = scanners.ScanWords
	}
	scanner.Split(splitFn)
	lineNumber := 1
	log.Debugln("start Scanner...")
	for scanner.Scan() {
		var matched bool
		text := scanner.Bytes()
		for _, token := range tokens {
			var err error
			var regexpPrefix string
			if *argInsensitive {
				regexpPrefix = "(?i)"
			}

			//fmt.Println("text ", text)
			//fmt.Println("token ", token)
			matched, err = regexp.MatchString(regexpPrefix+token, string(text))
			if err != nil {
				log.Println(err)
			}
			if matched {
				color := tmap[token]
				text = []byte(color.Sprint(string(text)))
				break
			}
		}
		if *argLineNumber {
			fmt.Printf("%d ", lineNumber)
			lineNumber++
		}
		if !*argOnlyMatching || matched && *argOnlyMatching {
			os.Stdout.Write(text)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	log.Debugf("end Scanner, processing %v completed.\n", path)
}
