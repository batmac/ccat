package highlighter

import (
	"ccat/log"
	"ccat/term"
	"ccat/utils"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

const (
	DEFAULT_STYLE = "monokai"
	//DEFAULT_FORMATTER = "terminal16m"
	MAX_READ_SIZE = 10_000_000
)

type Chroma struct {
	style     string
	formatter string
	lexer     string
}

func (h *Chroma) HighLight(w io.WriteCloser, r io.ReadCloser, o Options) error {
	log.Debugln(" highlighter: start chroma Highlighter")
	log.Debugln(log.Pp(o))

	var filename string = o.FileName

	// MAX_READ_SIZE Bytes max
	someSourceCode, err := ioutil.ReadAll(&io.LimitedReader{R: r, N: MAX_READ_SIZE})
	if err != nil {
		return err
	}

	_, err = r.Read(make([]byte, 1))
	if err != io.EOF {
		log.Fatal("highlighter: should read too much (file is too large for me)")
	}

	log.Debugf(" highlighter: read %v bytes\n", len(someSourceCode))
	if err := r.Close(); err != nil {
		log.Printf(" highlighter: %v\n", err)
	}

	// log.Debugf(" highlighter: registered lexers are: %v\n", lexers.Names(true))
	lexersList := lexers.Names(true)
	var lexer chroma.Lexer
	if len(o.LexerHint) > 0 && utils.StringInSlice(o.LexerHint, lexersList) {
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

	//log.Debugf(" highlighter: registered styles are: %v\n", styles.Names())
	//registered styles are: [abap algol algol_nu arduino autumn base16-snazzy borland bw colorful doom-one doom-one2 dracula emacs friendly fruity github hr_high_contrast hrdark igor lovelace manni monokai monokailight murphy native nord onesenterprise paraiso-dark paraiso-light pastie perldoc pygments rainbow_dash rrt solarized-dark solarized-dark256 solarized-light swapoff tango trac vim vs vulcan witchhazel xcode xcode-dark]

	stylesList := styles.Names()
	if o.StyleHint == "random" {
		rand.Seed(time.Now().UnixNano())
		randStyle := rand.Intn(len(stylesList))
		h.style = stylesList[randStyle]
	} else if len(o.StyleHint) > 0 && utils.StringInSlice(o.StyleHint, stylesList) {
		h.style = o.StyleHint
	} else {
		h.style = DEFAULT_STYLE
	}

	style := styles.Get(h.style) // or monokai
	if style == nil {
		style = styles.Fallback
	}
	log.Debugf(" highlighter: style is %+v\n", style.Name)

	log.Debugf(" highlighter: registered formatters are: %v\n", formatters.Names())

	formattersList := formatters.Names()
	if len(o.FormatterHint) > 0 && utils.StringInSlice(o.FormatterHint, formattersList) {
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (h *Chroma) help() string {
	return fmt.Sprintf("  - Lexers: %v\n  - Styles: %v\n  - Formatters: %v\n",
		lexers.Names(true),
		styles.Names(),
		formatters.Names(),
	)
}
