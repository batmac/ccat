package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/batmac/ccat/pkg/globalctx"
	"github.com/batmac/ccat/pkg/highlighter"
	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/mutators/pipeline"
	_ "github.com/batmac/ccat/pkg/mutators/simple"
	"github.com/batmac/ccat/pkg/openers"
	"github.com/batmac/ccat/pkg/pipedcmd"
	"github.com/batmac/ccat/pkg/scanners"
	"github.com/batmac/ccat/pkg/term"
	//_ "net/http/pprof"
)

func processFile(w io.Writer, path string) {
	log.Debugf("processing %s...\n", path)
	globalctx.Set("path", path)

	from, err := openers.Open(path, *argLockIn)
	if err != nil {
		log.Printf("opening %s: %v", path, err)
		setErrored()
		return
	}
	if len(*argMutators) > 0 {
		r, w := io.Pipe()
		err := pipeline.NewPipeline(*argMutators, w, from)
		if err != nil {
			setErrored()
			log.Fatal(err)
		}

		from = r
	}
	fromOrig := from
	defer func() {
		// I don't want to determine if already closed, try to close it, it will fail if it is already closed
		err := fromOrig.Close()
		if err != nil {
			log.Debugln(err)
		}
		log.Debugf("final closed %s...\n", path)
	}()
	/*************************************/
	if len(*argExec) > 0 {
		log.Debugf("creating pipedcmd %v...\n", *argExec)
		cmd, err := pipedcmd.New(*argExec)
		// log.Debugf("%s", log.Pp(cmd))
		if err != nil {
			setErrored()
			log.Fatal(err)
		}
		defer func() {
			log.Debugf("waiting pipedcmd %v...\n", *argExec)
			if err := cmd.Wait(); err != nil {
				setErrored()
				log.Println(err)
			}
		}()

		log.Debugf("start pipedcmd %s\n", cmd)

		err = cmd.Start(from)
		if err != nil {
			setErrored()
			log.Println(err)
		}
		from = cmd.Stdout.(*os.File)
	}

	if *argHuman {
		if term.IsArt(path) && len(*argMutators) == 0 {
			log.Debugf("is art, displaying directly...\n")
			term.PrintArt(from)
			return
		}
		expectingBinary := globalctx.Get("expectingBinary")
		if expectingBinary == nil || expectingBinary != nil && !expectingBinary.(bool) {
			log.Debugln("highlighting...")
			hl := globalctx.Get("hintLexer")
			if len(*argLexer) == 0 && hl != nil && len(hl.(string)) != 0 {
				hl := hl.(string)
				argLexer = &hl
			}

			r, w := io.Pipe()
			err := highlighter.Go(w, from, *highlighter.NewOptions(
				path,
				strings.ToLower(*argStyle),
				strings.ToLower(*argFormatter),
				*argLexer,
			))
			if err != nil {
				setErrored()
				log.Printf("error while highlighting: %v", err)
			} else {
				from = r
			}
		}
	}
	log.Debugln("initializing Scanner...")

	// go http.ListenAndServe(":8090", nil)
	scanner := bufio.NewScanner(from)

	splitFn := scanners.ScanBytes
	// splitFn := scanners.ScanLines
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

			// fmt.Println("text ", text)
			// fmt.Println("token ", token)
			matched, err = regexp.MatchString(regexpPrefix+token, string(text))
			if err != nil {
				setErrored()
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
			_, _ = w.Write(text)
		}
	}
	if err := scanner.Err(); err != nil {
		setErrored()
		log.Println(err)
	}
	log.Debugf("end Scanner, processing %v completed.\n", path)
	if len(*argMutators) > 0 {
		log.Debugln("Wait()ing pipeline...")
		pipeline.Wait()
	}
}

func setErrored() {
	globalctx.SetErrored()
}
