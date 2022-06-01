package main

import (
	"bufio"
	"ccat/pipedcmd"
	"ccat/scanners"
	"fmt"
	"log"
	"os"
	"regexp"
)

func processFile(path string) {
	from := os.Stdin
	var err error
	if path != "-" {
		from, err = fileOpen(path, *argLockIn)
		if err != nil {
			log.Println(err)
			return
		}
		defer fileClose(from, *argLockIn)
	}
	/*************************************/
	if len(*argExec) > 0 {
		log.Printf("creating pipedcmd %v\n", *argExec)
		cmd, err := pipedcmd.New(*argExec)
		if err != nil {
			log.Panicln(err)
		}
		defer func() {
			if err := cmd.Wait(); err != nil {
				log.Print(err)
			}
		}()

		print("start\n")

		err = cmd.Start(from)
		if err != nil {
			log.Panicln(err)
		}
		from = cmd.Stdout.(*os.File)
	}

	scanner := bufio.NewScanner(from)

	splitFn := scanners.ScanLines
	if *argSplitByWords {
		splitFn = scanners.ScanWords
	}
	scanner.Split(splitFn)
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
			fmt.Print(text)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
