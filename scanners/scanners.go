package scanners

import (
	"bytes"
	"unicode/utf8"
)

// All functions below are taken from go 1.18.2 and simplified because we don't want trimming.
// https://cs.opensource.google/go/go/+/refs/tags/go1.18.2:src/bufio/scan.go;l=351

func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[:i+1], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

// IsSpace reports whether the character is a Unicode white space character.
// We avoid dependency on the unicode package, but check validity of the implementation
// in the tests.
func IsSpace(r rune) bool {
	if r <= '\u00FF' {
		// Obvious ASCII ones: \t through \r plus space. Plus two Latin-1 oddballs.
		switch r {
		case ' ', '\t', '\n', '\v', '\f', '\r':
			return true
		case '\u0085', '\u00A0':
			return true
		}
		return false
	}
	// High-valued ones.
	if '\u2000' <= r && r <= '\u200a' {
		return true
	}
	switch r {
	case '\u1680', '\u2028', '\u2029', '\u202f', '\u205f', '\u3000':
		return true
	}
	return false
}

func ScanWords(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Scan until space, marking end of word.
	for width, i := 0, 0; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if IsSpace(r) {
			return i + width, data[:i+width], nil
		}
	}
	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > 0 {
		return len(data), data[:], nil
	}
	// Request more data.
	return 0, nil, nil
}

// ScanBytes is a split function for a Scanner that returns ALL data []bytes as a token.
func ScanBytes(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) < 4096 && !atEOF {
		// if len(data) < bufio.MaxScanTokenSize && !atEOF {
		return 0, nil, nil
	}
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	return len(data), data, nil
}
