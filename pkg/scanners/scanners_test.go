package scanners_test

import (
	//. "bufio"
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/batmac/ccat/pkg/scanners"
)

// const smallMaxTokenSize = 256 // Much smaller for more efficient testing.

// Test white space table matches the Unicode definition.
func TestSpace(t *testing.T) {
	for r := rune(0); r <= utf8.MaxRune; r++ {
		if scanners.IsSpace(r) != unicode.IsSpace(r) {
			t.Fatalf("white space property disagrees: %#U should be %t", r, unicode.IsSpace(r))
		}
	}
}

// slowReader is a reader that returns only a few bytes at a time, to test the incremental
// reads in Scanner.Scan.
type slowReader struct {
	max int
	buf io.Reader
}

func (sr *slowReader) Read(p []byte) (int, error) {
	if len(p) > sr.max {
		p = p[0:sr.max]
	}
	return sr.buf.Read(p)
}

// Test that the line splitter handles a final line without a newline.
func testNoNewline(t *testing.T, text string, lines []string) {
	t.Helper()
	buf := strings.NewReader(text)
	s := bufio.NewScanner(&slowReader{7, buf})
	s.Split(scanners.ScanLines)
	for lineNum := 0; s.Scan(); lineNum++ {
		line := lines[lineNum]
		if s.Text() != line {
			t.Errorf("%d: bad line: %d %d\n%.100q\n%.100q\n", lineNum, len(s.Bytes()), len(line), s.Bytes(), line)
		}
	}
	if err := s.Err(); err != nil {
		t.Fatal(err)
	}
}

// Test that the line splitter handles a final line without a newline.
func TestScanLineNoNewline(t *testing.T) {
	const text = "abcdefghijklmn\nopqrstuvwxyz"
	lines := []string{
		"abcdefghijklmn\n",
		"opqrstuvwxyz",
	}
	testNoNewline(t, text, lines)
}

// Test that the line splitter handles a final line with a carriage return but no newline.
func TestScanLineReturn(t *testing.T) {
	const text = "abcdefghijklmn\nopqrstuvwxyz\r"
	lines := []string{
		"abcdefghijklmn\n",
		"opqrstuvwxyz\r",
	}
	testNoNewline(t, text, lines)
}

// Test for issue 5268.
type alwaysError struct{}

func (alwaysError) Read(p []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func TestNonEOFWithEmptyRead(t *testing.T) {
	scanner := bufio.NewScanner(alwaysError{})
	scanner.Split(scanners.ScanLines)
	for scanner.Scan() {
		t.Fatal("read should fail")
	}
	if err := scanner.Err(); err != io.ErrUnexpectedEOF {
		t.Errorf("unexpected error: %v", err)
	}
}

// Test that Scan finishes if we have endless empty reads.
type endlessZeros struct{}

func (endlessZeros) Read(p []byte) (int, error) {
	return 0, nil
}

func TestBadReader(t *testing.T) {
	scanner := bufio.NewScanner(endlessZeros{})
	scanner.Split(scanners.ScanLines)
	for scanner.Scan() {
		t.Fatal("read should fail")
	}
	if err := scanner.Err(); err != io.ErrNoProgress {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestBlankLines(t *testing.T) {
	s := bufio.NewScanner(strings.NewReader(strings.Repeat("\n", 1000)))
	s.Split(scanners.ScanLines)
	for count := 0; s.Scan(); count++ {
		if count > 2000 {
			t.Fatal("looping")
		}
	}
	if s.Err() != nil {
		t.Fatal("after scan:", s.Err())
	}
}

func testScanWords(t *testing.T, text string, tokens []string) {
	t.Helper()
	buf := strings.NewReader(text)
	s := bufio.NewScanner(&slowReader{7, buf})
	s.Split(scanners.ScanWords)
	for tokNum := 0; s.Scan(); tokNum++ {
		token := tokens[tokNum]
		if s.Text() != token {
			t.Errorf("%d: bad token: %d %d\n%.100q\n%.100q\n", tokNum, len(s.Bytes()), len(token), s.Bytes(), token)
		}
	}
	if err := s.Err(); err != nil {
		t.Fatal(err)
	}
}

func TestScanWords(t *testing.T) {
	const text = "a beautiful day"
	words := []string{
		"a ",
		"beautiful ",
		"day",
	}
	testScanWords(t, text, words)
}

func testScanBytes(t *testing.T, text []byte, tokens []byte) {
	t.Helper()
	buf := bytes.NewReader(text)
	s := bufio.NewScanner(&slowReader{7, buf})
	s.Split(scanners.ScanBytes)
	for tokNum := 0; s.Scan(); tokNum++ {
		token := tokens[tokNum]
		if s.Bytes()[0] != token {
			t.Errorf("%d: bad token: %d %d\n%.100q\n%.100q\n", tokNum, len(s.Bytes()), 1, s.Bytes(), token)
		}
	}
	if err := s.Err(); err != nil {
		t.Fatal(err)
	}
}

func TestScanBytes(t *testing.T) {
	text := []byte("a beautiful day")
	bytes := []byte{'a', ' ', 'b', 'e', 'a', 'u', 't', 'i', 'f', 'u', 'l', ' ', 'd', 'a', 'y'}
	testScanBytes(t, text, bytes)
}
