package term_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/batmac/ccat/pkg/term"
)

func TestIsArt(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty", args{""}, false},
		{"exe", args{"xxx.exe"}, false},
		{"gif", args{"xxx.gif"}, true},
		{"png", args{"xxx.png"}, true},
		{"jpg", args{"xxx.jpg"}, true},
		{"jpeg", args{"xxx.jpeg"}, true},
		{"tiff", args{"xxx.tiff"}, true},
		{"tif", args{"xxx.tif"}, true},
		{"bmp", args{"xxx.bmp"}, true},
		{"webp", args{"xxx.webp"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := term.IsArt(tt.args.path); got != tt.want {
				t.Errorf("IsArt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrintArt(t *testing.T) {
	r, err := os.Open("testdata/blank.gif")
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		r io.ReadSeekCloser
	}
	tests := []struct {
		name string
		args args
	}{
		{"gif", args{r}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			term.PrintArt(tt.args.r)
		})
	}
}

func TestPrintANSIArt(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := term.PrintANSIArt(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("PrintANSIArt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	t.Run("panicplease", func(t *testing.T) {
		assertPanic(t, func() { _ = term.PrintANSIArt(strings.NewReader("not art")) })
	})
}

func assertPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}
