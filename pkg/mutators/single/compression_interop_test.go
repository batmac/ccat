package mutators_test

import (
	"archive/zip"
	"bytes"
	stdflate "compress/flate"
	stdgzip "compress/gzip"
	stdzlib "compress/zlib"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/batmac/ccat/pkg/globalctx"
	"github.com/batmac/ccat/pkg/mutators"
	dsbzip2 "github.com/dsnet/compress/bzip2"
	kzip "github.com/klauspost/compress/zip"
	"github.com/klauspost/compress/zstd"
)

// bigInput is compressible but not trivially so, and large enough to span
// several blocks of the parallel codecs (lz4: 4MiB, bzip2: 900kB).
func bigInput() string {
	var b strings.Builder
	for i := range 400_000 {
		fmt.Fprintf(&b, "line %d: %x\n", i, i*2654435761)
	}
	return b.String()
}

func mustClose(t *testing.T, c io.Closer) {
	t.Helper()
	if err := c.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestCompressionMultiBlockRoundTrip(t *testing.T) {
	globalctx.Set("path", "big.txt") // zip names its entry after it
	input := bigInput()
	for _, alg := range []string{"lz4", "zstd", "bzip2", "gzip", "zlib", "deflate", "zip"} {
		t.Run(alg, func(t *testing.T) {
			if got := mutators.Run("un"+alg, mutators.Run(alg, input)); got != input {
				t.Errorf("%s: round trip mismatch (len %d, want %d)", alg, len(got), len(input))
			}
		})
	}
}

// stdlib-produced streams must decode with ccat's decompressors.
func TestDecompressStdlibStreams(t *testing.T) {
	input := bigInput()
	var gz, zl, fl bytes.Buffer

	gw := stdgzip.NewWriter(&gz)
	_, _ = gw.Write([]byte(input))
	mustClose(t, gw)
	zw := stdzlib.NewWriter(&zl)
	_, _ = zw.Write([]byte(input))
	mustClose(t, zw)
	fw, _ := stdflate.NewWriter(&fl, stdflate.DefaultCompression)
	_, _ = fw.Write([]byte(input))
	mustClose(t, fw)

	tests := []struct{ mutator, data string }{
		{"ungzip", gz.String()},
		{"unzlib", zl.String()},
		{"undeflate", fl.String()},
		{"inflate", fl.String()},
	}
	for _, tt := range tests {
		t.Run(tt.mutator, func(t *testing.T) {
			if got := mutators.Run(tt.mutator, tt.data); got != input {
				t.Errorf("%s: mismatch (len %d, want %d)", tt.mutator, len(got), len(input))
			}
		})
	}
}

// ccat-produced streams must decode with the stdlib.
func TestCompressReadableByStdlib(t *testing.T) {
	input := bigInput()
	tests := []struct {
		mutator string
		reader  func(io.Reader) (io.Reader, error)
	}{
		{"gzip", func(r io.Reader) (io.Reader, error) { return stdgzip.NewReader(r) }},
		{"zlib", func(r io.Reader) (io.Reader, error) { return stdzlib.NewReader(r) }},
		{"deflate", func(r io.Reader) (io.Reader, error) { return stdflate.NewReader(r), nil }},
	}
	for _, tt := range tests {
		t.Run(tt.mutator, func(t *testing.T) {
			r, err := tt.reader(strings.NewReader(mutators.Run(tt.mutator, input)))
			if err != nil {
				t.Fatal(err)
			}
			got, err := io.ReadAll(r)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != input {
				t.Errorf("%s: mismatch (len %d, want %d)", tt.mutator, len(got), len(input))
			}
		})
	}
}

func TestDecompressConcatenatedStreams(t *testing.T) {
	bz := func(s string) string {
		var b bytes.Buffer
		w, err := dsbzip2.NewWriter(&b, nil)
		if err != nil {
			t.Fatal(err)
		}
		_, _ = w.Write([]byte(s))
		mustClose(t, w)
		return b.String()
	}
	gz := func(s string) string {
		var b bytes.Buffer
		w := stdgzip.NewWriter(&b)
		_, _ = w.Write([]byte(s))
		mustClose(t, w)
		return b.String()
	}

	if got := mutators.Run("unbzip2", bz("hello ")+bz("world")); got != "hello world" {
		t.Errorf("unbzip2: got %q", got)
	}
	if got := mutators.Run("ungzip", gz("hello ")+gz("world")); got != "hello world" {
		t.Errorf("ungzip: got %q", got)
	}
}

func TestUnzipMethods(t *testing.T) {
	const content = "hello from a zip archive, hello from a zip archive"

	t.Run("stdlib deflate", func(t *testing.T) {
		var b bytes.Buffer
		w := zip.NewWriter(&b)
		f, err := w.Create("a.txt")
		if err != nil {
			t.Fatal(err)
		}
		_, _ = f.Write([]byte(content))
		mustClose(t, w)
		if got := mutators.Run("unzip", b.String()); got != content {
			t.Errorf("got %q", got)
		}
	})

	for _, method := range []uint16{zstd.ZipMethodWinZip, zstd.ZipMethodPKWare} {
		t.Run(fmt.Sprintf("zstd method %d", method), func(t *testing.T) {
			var b bytes.Buffer
			w := kzip.NewWriter(&b)
			w.RegisterCompressor(method, zstd.ZipCompressor())
			f, err := w.CreateHeader(&kzip.FileHeader{Name: "a.txt", Method: method})
			if err != nil {
				t.Fatal(err)
			}
			_, _ = f.Write([]byte(content))
			mustClose(t, w)
			if got := mutators.Run("unzip", b.String()); got != content {
				t.Errorf("got %q", got)
			}
		})
	}
}
