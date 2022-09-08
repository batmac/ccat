package main

import (
	"bytes"
	"os"
	"reflect"
	"runtime"
	"testing"
)

func Test_printLicense(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipped on windows because crlf stuff.")
	}
	l, err := os.ReadFile("../../LICENSE")
	if err != nil {
		t.Fatal(err)
	}
	w := &bytes.Buffer{}
	printLicense(w)
	if gotW := w.Bytes(); !reflect.DeepEqual(gotW, l) {
		t.Errorf("printLicense() = %v, want %v", gotW, l)
	}
}
