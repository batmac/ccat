package main

import "strconv"

type ColorANSIbg int

const (
	Blackbg ColorANSIbg = iota
	Redbg
	Greenbg
	Yellowbg
	Bluebg
	Magentabg
	Cyanbg
	Whitebg
)

const Resetbg = "\x1b[0m"

var cachedbg = make(map[ColorANSIbg]string)

func (c ColorANSIbg) String() string {
	if _, ok := cachedbg[c]; !ok {
		cachedbg[c] = "\x1b[" + strconv.Itoa(40+int(c)) + "m"
	}
	return cachedbg[c]
}

func (c ColorANSIbg) Sprint(s string) string {
	return (c.String() + s + Resetbg)
}

func (c *ColorANSIbg) Next() Color {
	var n ColorANSIbg
	if *c == Cyanbg {
		n = Blackbg
	} else {
		n = *c + 1
	}
	return &n
}
