package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/batmac/ccat/color"
	"github.com/batmac/ccat/globalctx"
	"github.com/batmac/ccat/highlighter"
	"github.com/batmac/ccat/log"
	"github.com/batmac/ccat/mutators"
	"github.com/batmac/ccat/openers"
	"github.com/batmac/ccat/term"

	flag "github.com/spf13/pflag"
)

//go:generate go run gen.go

var (
	argTokens       = flag.StringP("tokens", "t", "", "comma-separated list of tokens")
	argInsensitive  = flag.BoolP("ignore-case", "i", false, "tokens given with -t are case-insensitive")
	argOnlyMatching = flag.BoolP("only", "o", false, "don't display lines without at least one token")
	argRaw          = flag.BoolP("raw", "r", false, "don't treat tokens as regexps")
	argLineNumber   = flag.BoolP("line-number", "n", false, "number the output lines, starting at 1.")
	argLockIn       = flag.BoolP("flock-in", "L", false, "exclusively flock each file before reading")
	argLockOut      = flag.BoolP("flock-out", "l", false, "exclusively flock stdout")
	argSplitByWords = flag.BoolP("word", "w", false, "read word by word instead of line by line (only works with utf8)")
	argExec         = flag.StringP("exec", "X", "", "command to exec on each file before processing it")
	argBG           = flag.BoolP("bg", "b", false, "colorize the background instead of the font")
	argDebug        = flag.BoolP("debug", "d", false, "debug what we are doing")
	argHuman        = flag.BoolP("humanize", "H", false, "try to do what is needed to help (syntax-highlight, autodetect, etc. TODO)")
	argStyle        = flag.StringP("style", "S", "", "style to use (only used if -H, look below for the list)")
	argFormatter    = flag.StringP("formatter", "F", "", "formatter to use (only used if -H, look below for the list)")
	argLexer        = flag.StringP("lexer", "P", "", "lexer to use (only used if -H, look below for the list)")
	argMutators     = flag.StringP("mutators", "m", "", "mutators to use (comma-separated)")
	argVersion      = flag.BoolP("version", "V", false, "print version on stdout")
	argLicense      = flag.Bool("license", false, "print license on stdout")
	argHelp         = flag.BoolP("help", "h", false, "print usage")
	argSelfUpdate   = flag.Bool("selfupdate", false, "Update to latest Github release")
	argCheckUpdate  = flag.Bool("check", false, "Check version with the latest Github release")

	tmap   map[string]color.Color
	tokens []string
)

var (
	// these are for the build tool
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
	builtBy = "unknown"
	tags    = "unknown"
)

func init() {
	flag.Usage = Usage
	flag.Parse()

	if !*argDebug {
		log.SetDebug(ioutil.Discard)
	} else {
		log.SetDebug(os.Stderr)
	}
}

func main() {
	log.Debugln("STARTING ccat")
	log.Debugf(buildLine())

	if *argVersion {
		fmt.Println(buildLine())
		os.Exit(0)
	}
	if *argLicense {
		fmt.Println(buildLine())
		printLicense()
		os.Exit(0)
	}
	if *argHelp {
		Usage()
		os.Exit(0)
	}
	if *argSelfUpdate {
		_ = update(version, false)
		os.Exit(0)
	}
	if *argCheckUpdate {
		_ = update(version, true)
		os.Exit(0)
	}
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

	// fmt.Printf("%v\n", tmap)

	log.Debugln("initializing file list...")
	fileList := flag.Args()
	if 0 == len(fileList) {
		fileList = []string{"-"}
	}
	log.Debugf("files: %v\n", fileList)

	setupStdout(*argLockOut)

	log.Debugln("processing...")
	for _, path := range fileList {
		globalctx.Reset()
		globalctx.Set("fileList", fileList)
		processFile(path)
	}

	if globalctx.IsErrored() {
		os.Exit(1)
	}
}

func buildLine() string {
	return fmt.Sprintf("version %s [%s], commit %s, built at %s by %s (%s)", version, tags, commit, date, builtBy, runtime.Version())
}

func Usage() {
	flag.CommandLine.SortFlags = false
	fmt.Fprintln(os.Stderr, buildLine())
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "---")

	fmt.Fprint(os.Stderr, "ccat <files>...\n")
	fmt.Fprint(os.Stderr, " - highlighter (used with -H):\n")
	fmt.Fprint(os.Stderr, highlighter.Help())
	fmt.Fprintf(os.Stderr, " - openers:\n    %v\n", strings.Join(openers.ListOpenersWithDescription(), "\n    "))
	fmt.Fprintf(os.Stderr, " - mutators:\n%v\n", availableMutatorsHelp())
}

func availableMutatorsHelp() string {
	var s strings.Builder
	l := mutators.ListAvailableMutatorsByCategoryWithDescriptions()
	keys := make([]string, 0, len(l))
	for k := range l {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, category := range keys {
		if len(category) > 0 {
			s.WriteString("    " + category + ":\n")
		}
		sort.Strings(l[category])
		for _, mutator := range l[category] {
			s.WriteString("        " + mutator + "\n")
		}
	}
	return s.String()
}
