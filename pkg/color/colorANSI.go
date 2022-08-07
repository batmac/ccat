package color

import "strconv"

type ANSI int

const (
	Black ANSI = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

const Reset = "\x1b[0m"

var cached = make(map[ANSI]string)

func (c ANSI) String() string {
	if _, ok := cached[c]; !ok {
		cached[c] = "\x1b[" + strconv.Itoa(30+int(c)) + "m"
	}
	return cached[c]
}

func (c ANSI) Sprint(s string) string {
	return (c.String() + s + Reset)
}

func (c *ANSI) Next() Color {
	var n ANSI
	if *c == Cyan {
		n = Black
	} else {
		n = *c + 1
	}
	return &n
}
