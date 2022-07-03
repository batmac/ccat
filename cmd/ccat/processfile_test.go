package main

import (
	"os"
	"testing"

	_ "github.com/batmac/ccat/pkg/mutators/simple"
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

	discardFile, err := os.OpenFile(os.DevNull, os.O_APPEND, 0)
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = discardFile
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processFile(tt.args.path)
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
		t.Run(tt.name, func(t *testing.T) {
			setError()
		})
	}
}
