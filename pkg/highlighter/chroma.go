//go:build !nohl
// +build !nohl

package highlighter

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/stringutils"
	"github.com/batmac/ccat/pkg/term"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

const (
	DefaultStyle = "monokai"
	// DEFAULT_FORMATTER = "terminal16m"
	MaxReadSize = 10_000_000
)

type Chroma struct {
	style     string
	formatter string
	lexer     string
}

func (h *Chroma) highLight(w io.WriteCloser, r io.ReadCloser, o Options) error {
	log.Debugln(" highlighter: start chroma Highlighter")
	// log.Debugln(log.Pp(o))

	filename := o.FileName

	// MAX_READ_SIZE Bytes max
	someSourceCode, err := io.ReadAll(&io.LimitedReader{R: r, N: MaxReadSize})
	if err != nil {
		return err
	}

	additionalChar := make([]byte, 1)
	if _, err := r.Read(additionalChar); err == nil {
		log.Debugf("highlighter: output is too large for me (> %s), I will not highlight it", stringutils.HumanSize(MaxReadSize))
		_, err = w.Write(someSourceCode)
		if err != nil {
			log.Printf(" highlighter: %v\n", err)
		}
		_, err = w.Write(additionalChar)
		if err != nil {
			log.Printf(" highlighter: %v\n", err)
		}
		_, err = io.Copy(w, r)
		if err != nil {
			log.Printf(" highlighter: %v\n", err)
		}
		return nil
	}

	log.Debugf(" highlighter: read %v bytes\n", len(someSourceCode))
	if err := r.Close(); err != nil {
		log.Printf(" highlighter: %v\n", err)
	}

	// log.Debugf(" highlighter: registered lexers are: %v\n", lexers.Names(true))
	lexersList := lexers.Names(true)
	var lexer chroma.Lexer
	if checkWithFuzzy(o.LexerHint, lexersList) {
		log.Debugf(" highlighter: setting the lexer to %v\n", o.LexerHint)
		lexer = lexers.Get(o.LexerHint)
	} else {
		lexer = lexers.Match(filename)
		if lexer == nil {
			log.Debugf(" highlighter: filename did not help to find a lexer, analyzing content...\n")
			lexer = lexers.Analyse(string(someSourceCode))
			if lexer == nil {
				log.Debugf(" highlighter: fallbacking the lexer\n")
				lexer = lexers.Fallback
			}
		}
		lexer = chroma.Coalesce(lexer)
	}

	log.Debugf(" highlighter: chosen Lexer is %v\n", lexer.Config().Name)
	h.lexer = lexer.Config().Name

	// log.Debugf(" highlighter: registered styles are: %v\n", styles.Names())
	// registered styles are: [abap algol algol_nu arduino autumn base16-snazzy borland bw colorful doom-one doom-one2 dracula emacs friendly fruity github hr_high_contrast hrdark igor lovelace manni monokai monokailight murphy native nord onesenterprise paraiso-dark paraiso-light pastie perldoc pygments rainbow_dash rrt solarized-dark solarized-dark256 solarized-light swapoff tango trac vim vs vulcan witchhazel xcode xcode-dark]

	stylesList := styles.Names()
	switch {
	case o.StyleHint == "random":
		rand.Seed(time.Now().UnixNano())
		//#nosec G404 (weak rand)
		randStyle := rand.Intn(len(stylesList))
		h.style = stylesList[randStyle]
	case checkWithFuzzy(o.StyleHint, stylesList):
		h.style = o.StyleHint
	default:
		h.style = DefaultStyle
	}

	style := styles.Get(h.style)
	if style == nil {
		style = styles.Fallback
	}
	log.Debugf(" highlighter: style is %+v\n", style.Name)

	log.Debugf(" highlighter: registered formatters are: %v\n", formatters.Names())

	formattersList := formatters.Names()
	if checkWithFuzzy(o.FormatterHint, formattersList) {
		h.formatter = o.FormatterHint
	} else {
		c := term.SupportedColors()
		switch {
		case c >= 16_000_000:
			h.formatter = "terminal16m"
		case c >= 256:
			h.formatter = "terminal256"
		case c >= 16:
			h.formatter = "terminal16"
		case c >= 8:
			h.formatter = "terminal8"
		default:
			h.formatter = "noop"
		}
		/* 		h.formatter = DEFAULT_FORMATTER */
	}
	formatter := formatters.Get(h.formatter)
	if formatter == nil {
		formatter = formatters.Fallback
	}
	log.Debugf(" highlighter: formatter is %v\n", h.formatter)

	iterator, err := lexer.Tokenise(nil, string(someSourceCode))
	if err != nil {
		return err
	}

	err = formatter.Format(w, style, iterator)
	if err != nil {
		return err
	}

	log.Debugln(" highlighter: end chroma Highlight")
	return nil
}

func checkWithFuzzy(s string, list []string) bool {
	if len(s) == 0 {
		return false
	}
	// log.Printf("%v\n", list)
	if stringutils.IsStringInSlice(s, list) {
		return true
	}
	fs, err := stringutils.FuzzySearch(s, list, 0.5)
	if err != nil {
		log.Fatal(err)
	}
	if len(fs) == 0 {
		return false
	}
	fmt.Fprintf(os.Stderr, "'%s' does not exist, did you mean %s?\n", s, fs)
	os.Exit(1)
	return false
}

func (h *Chroma) help() string {
	lexers, styles, formatters := lexers.Names(true), styles.Names(), formatters.Names()
	stringutils.SortStringsCaseInsensitive(lexers)
	stringutils.SortStringsCaseInsensitive(styles)
	stringutils.SortStringsCaseInsensitive(formatters)
	return fmt.Sprintf("  - Lexers: %v\n  - Styles: %v\n  - Formatters: %v\n",
		strings.Join(lexers, ", "),
		strings.Join(styles, ", "),
		strings.Join(formatters, ", "),
	)
}
