package main

import (
	"io"
	"os"
	"testing"

	_ "github.com/batmac/ccat/pkg/mutators/single"
)

func Test_processFile(t *testing.T) {
	exe, err := os.Executable() // we just need an existing file
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
	}{
		{"donotpanicplease on not existing", args{"not existing"}},
		{"donotpanicplease on dir", args{"."}},
		{"donotpanicplease on exe", args{exe}},
		{"donotpanicplease on exe", args{"file://" + exe}},
		{"donotpanicplease on exe", args{"fake://fakefakefake"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			processFile(io.Discard, tt.args.path)
			processFileAsIs(io.Discard, tt.args.path)
		})
	}
}

func Test_setError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"donotpanicplease"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			setErrored()
		})
	}
}
