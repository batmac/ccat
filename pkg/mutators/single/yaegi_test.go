//go:build plugins
// +build plugins

package mutators_test

import (
	"os"
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
	"github.com/batmac/ccat/pkg/stringutils"
)

var code = `
package truc

import (
	"fmt"
	"io"
	"os"
)

func Y(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	count, err := io.Copy(io.Discard, r)
	if err != nil {
		return 0, err
	}
	fmt.Fprintln(w, count)
	return -1, nil
}
`

func Test_applyYaegi(t *testing.T) {
	tests := []struct {
		name, input, code, symbol, output string
	}{
		{
			name:   "abc",
			input:  "abc",
			code:   code,
			symbol: "truc.Y",
			output: "3",
		},
	}

	f := "yaegi"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// write "tt.code" in a tmp file
			file, err := os.CreateTemp("", "ccat-yaegi-")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(file.Name())
			if _, err := file.WriteString(tt.code); err != nil {
				t.Fatal(err)
			}
			if err := file.Close(); err != nil {
				t.Fatal(err)
			}

			if got := mutators.Run(f+":"+file.Name()+":"+tt.symbol, tt.input); stringutils.DeleteSpaces(got) != stringutils.DeleteSpaces(tt.output) {
				t.Errorf("%s = %v, want %v", f, got, tt.output)
			}
		})
	}
}
