package main

import "strconv"

type ColorANSI int

const (
	Black ColorANSI = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

const Reset = "\x1b[0m"

var cached = make(map[ColorANSI]string)

func (c ColorANSI) String() string {
	if _, ok := cached[c]; !ok {
		cached[c] = "\x1b[" + strconv.Itoa(30+int(c)) + "m"
	}
	return cached[c]
}

func (c ColorANSI) Sprint(s string) string {
	return (c.String() + s + Reset)
}

func (c *ColorANSI) Next() Color {
	var n ColorANSI
	if *c == Cyan {
		n = Black
	} else {
		n = *c + 1
	}
	return &n
}
