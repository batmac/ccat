package color

import (
	"math/rand"
	"strconv"
)

type Color256 int

const CReset = "\x1b[0m"

var ccached = make(map[Color256]string)

func (c Color256) String() string {
	if _, ok := ccached[c]; !ok {
		ccached[c] = "\x1b[38;5;" + strconv.Itoa(1+int(c)) + "m"
	}
	return ccached[c]
}

func (c Color256) Sprint(s string) string {
	return c.String() + s + CReset
}

//#nosec
func (c *Color256) Next() Color {
	rand.Seed(int64(*c))
	n := Color256(rand.Intn(230))
	return &n
}
