package term_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/batmac/ccat/term"
)

func TestIsITerm2(t *testing.T) {
	t.Run("donotpanicplease", func(t *testing.T) {
		_ = term.IsITerm2()
	})
}

func TestPrintITerm2Art(t *testing.T) {
	r, err := os.Open("testdata/blank.gif")
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"gif", args{r}, false},
		{"not art", args{strings.NewReader("not art")}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := term.PrintITerm2Art(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("PrintITerm2Art() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
