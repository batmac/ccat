package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/batmac/ccat/pkg/color"
	"github.com/batmac/ccat/pkg/globalctx"
	"github.com/batmac/ccat/pkg/highlighter"
	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/mutators"
	"github.com/batmac/ccat/pkg/openers"
	"github.com/batmac/ccat/pkg/term"

	flag "github.com/spf13/pflag"
)

//go:generate go run gen_licence.go
//go:generate go mod tidy
//go:generate go run gen_gomod.go

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
	argHuman        = flag.BoolP("humanize", "H", false, "try to do what is needed to help (syntax-highlight, autodetect, etc.)")
	argStyle        = flag.StringP("style", "S", "", "style to use (only used if -H, --fullhelp for the list)")
	argFormatter    = flag.StringP("formatter", "F", "", "formatter to use (only used if -H, --fullhelp for the list)")
	argLexer        = flag.StringP("lexer", "P", "", "lexer to use (only used if -H, --fullhelp for the list)")
	argMutators     = flag.StringP("mutators", "m", "", "mutators to use (comma-separated), --fullhelp for the list")
	argVersion      = flag.BoolP("version", "V", false, "print version on stdout")
	argLicense      = flag.Bool("license", false, "print license on stdout")
	argGomod        = flag.Bool("gomod", false, "print used go.mod on stdout")
	argHelp         = flag.BoolP("help", "h", false, "print usage")
	argFullHelp     = flag.BoolP("fullhelp", "", false, "print full usage")
	argSelfUpdate   = flag.Bool("selfupdate", false, "Update to latest Github release")
	argCheckUpdate  = flag.Bool("check", false, "Check version with the latest Github release")
	argDebug        = flag.BoolP("debug", "d", false, "debug what we are doing")
	argInsecure     = flag.BoolP("insecure", "k", false, "get files insecurely (globally)")

	tmap   map[string]color.Color
	tokens []string
)

var (
	// these are for the build tool
	version = "unknown"
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
		printLicense(os.Stdout)
		os.Exit(0)
	}
	if *argGomod {
		b := strings.Builder{}
		b.WriteString("# " + buildLine())
		bi, ok := debug.ReadBuildInfo()
		if !ok {
			log.Printf("failed to read build info")
		}
		b.WriteString("\n\n## Modules used:\n")
		for _, d := range bi.Deps {
			fmt.Fprintf(&b, "\n- %v %v (%v)", d.Path, d.Version, d.Sum)
		}
		var s string
		if *argHuman {
			s = highlighter.Run(b.String(), highlighter.NewOptions(
				"",
				strings.ToLower(*argStyle),
				strings.ToLower(*argFormatter),
				"go"))
		} else {
			s = b.String()
		}
		fmt.Println(s)
		os.Exit(0)
	}
	if *argHelp {
		Usage()
		os.Exit(0)
	}
	if *argFullHelp {
		FullUsage()
		os.Exit(0)
	}
	if *argSelfUpdate {
		update(version, false)
		os.Exit(0)
	}
	if *argCheckUpdate {
		update(version, true)
		os.Exit(0)
	}

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
	if len(fileList) == 0 {
		fileList = []string{"-"}
	}
	log.Debugf("files: %v\n", fileList)

	if err := setupStdout(*argLockOut); err != nil {
		log.Fatal(err)
	}

	log.Debugln("processing...")
	process := processFile
	if !term.IsStdoutTerminal() && (flag.NFlag() == 0 || (flag.NFlag() == 1 && *argDebug)) {
		log.Debugln("no option given, trying to be as fast as possible")
		process = processFileAsIs
	}
	for _, path := range fileList {
		globalctx.Reset()
		globalctx.Set("fileList", fileList)
		globalctx.Set("insecure", *argInsecure)

		process(path)
	}

	if globalctx.IsErrored() {
		os.Exit(1)
	}
}

func buildLine() string {
	return fmt.Sprintf("version %s [%s], commit %s, built at %s by %s (%s %s/%s)", version, tags, commit, date, builtBy, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

func Usage() {
	usage(false)
}

func FullUsage() {
	usage(true)
}

func usage(full bool) {
	flag.CommandLine.SortFlags = false
	fmt.Fprintln(os.Stderr, buildLine())
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "")

	if !full {
		return
	}

	fmt.Fprintln(os.Stderr, "---")
	fmt.Fprint(os.Stderr, "ccat <files>...\n")
	fmt.Fprint(os.Stderr, " - highlighter (used with -H):\n")
	fmt.Fprint(os.Stderr, highlighter.Help())
	fmt.Fprintf(os.Stderr, " - openers:\n    %v\n", strings.Join(openers.ListOpenersWithDescription(), "\n    "))
	fmt.Fprintf(os.Stderr, " - mutators:\n%v\n", mutators.AvailableMutatorsHelp())
}
