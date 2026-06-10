# go-gemini

[![godocs.io](https://godocs.io/git.sr.ht/~adnano/go-gemini?status.svg)](https://godocs.io/git.sr.ht/~adnano/go-gemini) [![builds.sr.ht status](https://builds.sr.ht/~adnano/go-gemini.svg)](https://builds.sr.ht/~adnano/go-gemini?)

Package gemini implements the [Gemini protocol](https://geminiprotocol.net)
in Go. It provides an API similar to that of net/http to facilitate the
development of Gemini clients and servers.

Compatible with version v0.16.0 of the Gemini specification.

## Usage

	import "git.sr.ht/~adnano/go-gemini"

Note that some filesystem-related functionality is only available on Go 1.16
or later as it relies on the io/fs package.

## Examples

There are a few examples provided in the examples directory.
To run an example:

	go run examples/server.go

## License

go-gemini is licensed under the terms of the MIT license (see LICENSE).
Portions of this library were adapted from Go and are governed by a BSD-style
license (see LICENSE-GO). Those files are marked accordingly.
