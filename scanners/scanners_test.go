package scanners_test

import (
	//. "bufio"
	"bufio"
	"io"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/batmac/ccat/scanners"
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

func (sr *slowReader) Read(p []byte) (n int, err error) {
	if len(p) > sr.max {
		p = p[0:sr.max]
	}
	return sr.buf.Read(p)
}

// Test that the line splitter handles a final line without a newline.
func testNoNewline(text string, lines []string, t *testing.T) {
	buf := strings.NewReader(text)
	s := bufio.NewScanner(&slowReader{7, buf})
	s.Split(scanners.ScanLines)
	for lineNum := 0; s.Scan(); lineNum++ {
		line := lines[lineNum]
		if s.Text() != line {
			t.Errorf("%d: bad line: %d %d\n%.100q\n%.100q\n", lineNum, len(s.Bytes()), len(line), s.Bytes(), line)
		}
	}
	err := s.Err()
	if err != nil {
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
	testNoNewline(text, lines, t)
}

// Test that the line splitter handles a final line with a carriage return but no newline.
func TestScanLineReturn(t *testing.T) {
	const text = "abcdefghijklmn\nopqrstuvwxyz\r"
	lines := []string{
		"abcdefghijklmn\n",
		"opqrstuvwxyz\r",
	}
	testNoNewline(text, lines, t)
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
	err := scanner.Err()
	if err != io.ErrUnexpectedEOF {
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
	err := scanner.Err()
	if err != io.ErrNoProgress {
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
