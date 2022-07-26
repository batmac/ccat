package highlighter

type Options struct {
	FileName      string
	StyleHint     string
	LexerHint     string
	FormatterHint string
}

func NewOptions(filename, style, formatter, lexer string) *Options {
	return &Options{
		FileName:      filename,
		StyleHint:     style,
		FormatterHint: formatter,
		LexerHint:     lexer,
	}
}
