# ccat [![Go](https://github.com/batmac/ccat/actions/workflows/go.yml/badge.svg)](https://github.com/batmac/ccat/actions/workflows/go.yml) ![GitHub](https://img.shields.io/github/license/batmac/ccat) [![Go Report Card](https://goreportcard.com/badge/github.com/batmac/ccat)](https://goreportcard.com/report/github.com/batmac/ccat)
cat on steroids


## build
you need go >=1.15, available build tags:
- `libcurl`: build with the libcurl opener.
- `fileonly`: build with the local file opener only.
- `nomd`: build without the markdown interpreter (glamour).
- `nohl`: build without the syntax-highlighter.
- `crappy`: build with some crappy (but useful) openers/mutators (needs a recent go version).

for instance:
`go build --tags libcurl,crappy .`

## help

