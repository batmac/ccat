package utils_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/batmac/ccat/utils"
)

var (
	taint       bool
	emptyReader = strings.NewReader("")
	hello       = "hello"
	helloReader = strings.NewReader(hello)
	ErrHello    = errors.New(hello)
)

func TestNewReadCloser(t *testing.T) {
	type args struct {
		r       io.Reader
		closure func() error
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		taint   bool
	}{
		{"empty", args{emptyReader, func() error { return nil }}, nil, false},
		{"simple", args{helloReader, func() error {
			taint = true
			return nil
		}}, nil, true},
		{"hello", args{helloReader, func() error {
			taint = true
			return ErrHello
		}}, ErrHello, true},
	}
	for _, tt := range tests {
		taint = false
		t.Run(tt.name, func(t *testing.T) {
			got := utils.NewReadCloser(tt.args.r, tt.args.closure)
			err := got.Close()
			if err != tt.wantErr || tt.taint != taint {
				t.Errorf("NewReadCloser() failed: err = %v, wantErr = %v, taint = %v, want %v", err, tt.wantErr, taint, tt.taint)
			}
		})
	}
}

func Test_newCloser_Close(t *testing.T) {
	type fields struct {
		Reader  io.Reader
		closure func() error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"empty", fields{emptyReader, func() error { return nil }}, false},
		{"simple", fields{helloReader, func() error { return fmt.Errorf(hello) }}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := utils.NewReadCloser(tt.fields.Reader, tt.fields.closure)

			if err := c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("newCloser.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newCloserWriterTo_Close(t *testing.T) {
	type fields struct {
		Reader  io.Reader
		closure func() error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"empty", fields{emptyReader, func() error { return nil }}, false},
		{"simple", fields{helloReader, func() error { return fmt.Errorf(hello) }}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := utils.NewReadCloser(tt.fields.Reader, tt.fields.closure)
			if err := c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("newCloserWriterTo.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newCloserWriterTo_WriteTo(t *testing.T) {
	type fields struct {
		Reader  io.Reader
		closure func() error
	}
	tests := []struct {
		name    string
		fields  fields
		wantN   int64
		wantW   string
		wantErr bool
	}{
		{"empty", fields{emptyReader, func() error { return nil }}, 0, "", false},
		{"simple", fields{helloReader, func() error { return nil }}, int64(len(hello)), "hello", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := utils.NewReadCloser(tt.fields.Reader, tt.fields.closure)

			w := &bytes.Buffer{}
			gotN, err := c.(io.WriterTo).WriteTo(w)
			if (err != nil) != tt.wantErr {
				t.Errorf("newCloserWriterTo.WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("newCloserWriterTo.WriteTo() = %v, want %v", gotN, tt.wantN)
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("newCloserWriterTo.WriteTo() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
