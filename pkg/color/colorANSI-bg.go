package color

import "strconv"

type ANSIbg int

const (
	Blackbg ANSIbg = iota
	Redbg
	Greenbg
	Yellowbg
	Bluebg
	Magentabg
	Cyanbg
	Whitebg
)

const Resetbg = "\x1b[0m"

var cachedbg = make(map[ANSIbg]string)

func (c ANSIbg) String() string {
	if _, ok := cachedbg[c]; !ok {
		cachedbg[c] = "\x1b[" + strconv.Itoa(40+int(c)) + "m"
	}
	return cachedbg[c]
}

func (c ANSIbg) Sprint(s string) string {
	return (c.String() + s + Resetbg)
}

func (c *ANSIbg) Next() Color {
	var n ANSIbg
	if *c == Cyanbg {
		n = Blackbg
	} else {
		n = *c + 1
	}
	return &n
}
