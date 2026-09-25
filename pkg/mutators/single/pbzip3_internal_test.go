package mutators

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/jftuga/bzip3-go"
)

type bufWriteCloser struct{ bytes.Buffer }

func (*bufWriteCloser) Close() error { return nil }

// failingWriter fails after accepting `left` bytes
type failingWriter struct{ left int }

var errTestWrite = errors.New("write failed")

func (f *failingWriter) Write(p []byte) (int, error) {
	if len(p) > f.left {
		return 0, errTestWrite
	}
	f.left -= len(p)
	return len(p), nil
}
func (*failingWriter) Close() error { return nil }

type failingReader struct{ r io.Reader }

var errTestRead = errors.New("read failed")

func (f failingReader) Read(p []byte) (int, error) {
	n, err := f.r.Read(p)
	if errors.Is(err, io.EOF) {
		return n, errTestRead
	}
	return n, err
}

var pbzip3TestInput = strings.Repeat("Lorem ipsum dolor sit amet, consectetur adipiscing elit. ", 8000)

func pbzip3Compressed(t *testing.T) string {
	t.Helper()
	out := new(bufWriteCloser)
	if _, err := cpbzip3(out, io.NopCloser(strings.NewReader(pbzip3TestInput)), pbzip3Config{workers: 4, blockSize: bzip3.MinBlockSize}); err != nil {
		t.Fatal(err)
	}
	return out.String()
}

func Test_punbzip3Errors(t *testing.T) {
	compressed := pbzip3Compressed(t)
	tests := map[string]string{
		"bad magic":        "XZ3v1" + compressed[5:],
		"short header":     compressed[:5],
		"bad block size":   "BZ3v1\x00\x00\x00\x00" + compressed[9:],
		"truncated header": compressed[:9+4],
		"truncated block":  compressed[:len(compressed)-10],
		"corrupted block":  compressed[:40] + "\xff\xff\xff\xff" + compressed[44:],
	}
	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := punbzip3(new(bufWriteCloser), io.NopCloser(strings.NewReader(input)), 4); err == nil {
				t.Errorf("%s: expected an error", name)
			}
		})
	}
}

func Test_pbzip3PropagatesErrors(t *testing.T) {
	conf := pbzip3Config{workers: 3, blockSize: bzip3.MinBlockSize}

	compressed := pbzip3Compressed(t)

	// fails on a later block while other blocks are in flight
	_, err := cpbzip3(&failingWriter{left: len(compressed) / 2}, io.NopCloser(strings.NewReader(pbzip3TestInput)), conf)
	if !errors.Is(err, errTestWrite) {
		t.Errorf("compress, write error: got %v", err)
	}

	_, err = cpbzip3(new(bufWriteCloser), io.NopCloser(failingReader{strings.NewReader(pbzip3TestInput)}), conf)
	if !errors.Is(err, errTestRead) {
		t.Errorf("compress, read error: got %v", err)
	}

	_, err = punbzip3(&failingWriter{left: len(pbzip3TestInput) / 2}, io.NopCloser(strings.NewReader(compressed)), 3)
	if !errors.Is(err, errTestWrite) {
		t.Errorf("decompress, write error: got %v", err)
	}
}

func Test_pbzip3ReturnedSizes(t *testing.T) {
	compressed := pbzip3Compressed(t)
	out := new(bufWriteCloser)
	n, err := punbzip3(out, io.NopCloser(strings.NewReader(compressed)), 2)
	if err != nil {
		t.Fatal(err)
	}
	if n != int64(len(pbzip3TestInput)) || out.String() != pbzip3TestInput {
		t.Errorf("punbzip3 returned %d, want %d", n, len(pbzip3TestInput))
	}

	n, err = cpbzip3(new(bufWriteCloser), io.NopCloser(strings.NewReader(pbzip3TestInput)), pbzip3Config{workers: 2, blockSize: bzip3.MinBlockSize})
	if err != nil || n != int64(len(pbzip3TestInput)) {
		t.Errorf("pbzip3 returned %d, %v, want %d", n, err, len(pbzip3TestInput))
	}
}
