# lzfse

[![Go](https://github.com/aixiansheng/lzfse/workflows/Go/badge.svg?branch=master)](https://github.com/aixiansheng/lzfse/actions) [![GoDoc](https://godoc.org/github.com/aixiansheng/lzfse?status.svg)](https://pkg.go.dev/github.com/aixiansheng/lzfse)

> An LZFSE decompressor written in Go

```golang
package main

import (
	"os"
	"gihub.com/aixiansheng/lzfse"
)

func main() {
	inf, err := os.Open("some.lzfse")
	outf, err := os.Create("some.file")
	d := lzfse.NewReader(fh)
	io.Copy(outf, d)
}
```

## Testing

```
make -C test/

# all tests
go test -v

# just one test
go test -v -run TestVariousSizes/test/test.small.dec.cmp
```
