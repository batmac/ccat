package color

import (
	"math/rand"
	"strconv"
)

type C256 int

const CReset = "\x1b[0m"

var ccached = make(map[C256]string)

func (c C256) String() string {
	if _, ok := ccached[c]; !ok {
		ccached[c] = "\x1b[38;5;" + strconv.Itoa(1+int(c)) + "m"
	}
	return ccached[c]
}

func (c C256) Sprint(s string) string {
	return c.String() + s + CReset
}

// #nosec
func (c *C256) Next() Color {
	rand.Seed(int64(*c))      //nolint:staticcheck
	n := C256(rand.Intn(230)) //nolint:gosec
	return &n
}
