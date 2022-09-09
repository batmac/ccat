package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/batmac/ccat/pkg/color"
	"github.com/batmac/ccat/pkg/completion"
	"github.com/batmac/ccat/pkg/globalctx"
	"github.com/batmac/ccat/pkg/highlighter"
	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/mutators"
	"github.com/batmac/ccat/pkg/openers"
	"github.com/batmac/ccat/pkg/selfupdate"
	"github.com/batmac/ccat/pkg/term"

	flag "github.com/spf13/pflag"
)

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
	argBuildInfo    = flag.BoolP("buildinfo", "B", false, "print build info on stdout")
	argHelp         = flag.BoolP("help", "h", false, "print usage on stderr")
	argFullHelp     = flag.BoolP("fullhelp", "", false, "print full usage on stdout")
	argSelfUpdate   = flag.Bool("selfupdate", false, "Update to latest Github release")
	argCheckUpdate  = flag.Bool("check", false, "Check version with the latest Github release")
	argDebug        = flag.BoolP("debug", "d", false, "debug what we are doing")
	argInsecure     = flag.BoolP("insecure", "k", false, "get files insecurely (globally)")
	argCompletion   = flag.StringP("completion", "C", "", "print shell completion script")
	argLess         = flag.BoolP("ui", "T", false, "display with a minimal ui")

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
		log.SetDebug(io.Discard)
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
	if *argBuildInfo {
		b := strings.Builder{}
		b.WriteString(buildLine())
		bi, ok := debug.ReadBuildInfo()
		if !ok {
			log.Fatal("failed to read build info")
		}
		b.WriteString(bi.String())

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
	if len(*argCompletion) > 0 {
		var opts []string
		flag.VisitAll(func(f *flag.Flag) {
			if len(f.Shorthand) > 0 {
				opts = append(opts, "-"+f.Shorthand)
			}
			opts = append(opts, "--"+f.Name)
		})
		completion.Print(*argCompletion, opts)
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
		selfupdate.Do(version, tags, false)
		os.Exit(0)
	}
	if *argCheckUpdate {
		selfupdate.Do(version, tags, true)
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
			c = new(color.ANSIbg)
		} else {
			c = new(color.ANSI)
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
	if *argLess {
		if *argDebug {
			log.Fatal("aborting, because the ui option is not compatible with the debug option")
		}
		process = uiWrapProcessFile(process)
	}
	for _, path := range fileList {
		globalctx.Reset()
		globalctx.Set("fileList", fileList)
		globalctx.Set("insecure", *argInsecure)

		process(os.Stdout, path)
	}

	if globalctx.IsErrored() {
		os.Exit(1)
	}
}

func buildLine() string {
	return fmt.Sprintf("version %s [%s], commit %s, built at %s by %s (%s %s/%s)", version, tags, commit, date, builtBy, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

func Usage() {
	usage(os.Stderr, false)
}

func FullUsage() {
	usage(os.Stdout, true)
}

func usage(w io.Writer, full bool) {
	flag.CommandLine.SortFlags = false
	flag.CommandLine.SetOutput(w)
	fmt.Fprintln(w, buildLine())
	flag.PrintDefaults()
	fmt.Fprintln(w, "")

	if !full {
		return
	}

	fmt.Fprintln(w, "---")
	fmt.Fprint(w, "ccat <files>...\n")
	fmt.Fprint(w, " - highlighter (used with -H):\n")
	fmt.Fprint(w, highlighter.Help())
	fmt.Fprintf(w, " - openers:\n    %v\n", strings.Join(openers.ListOpenersWithDescription(), "\n    "))
	fmt.Fprintf(w, " - mutators:\n%v\n", mutators.AvailableMutatorsHelp())
}
