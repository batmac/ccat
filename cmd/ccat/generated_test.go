package main

import (
	"bytes"
	"os"
	"reflect"
	"testing"
)

func Test_printLicense(t *testing.T) {
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
