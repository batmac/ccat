package highlighter_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/batmac/ccat/pkg/highlighter"
	"github.com/batmac/ccat/pkg/utils"
)

func TestGo(t *testing.T) {
	type args struct {
		w io.WriteCloser
		r io.ReadCloser
		o highlighter.Options
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{&utils.NopStringWriteCloser{}, ioutil.NopCloser(&bytes.Buffer{}), highlighter.Options{"", "", "", ""}}, false},
		{"simple", args{&utils.NopStringWriteCloser{}, ioutil.NopCloser(bytes.NewBufferString("hello")), highlighter.Options{"", "", "", ""}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := highlighter.Go(tt.args.w, tt.args.r, tt.args.o); (err != nil) != tt.wantErr {
				t.Errorf("Go() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHelp(t *testing.T) {
	t.Run("donotpanicplease", func(t *testing.T) {
		if got := highlighter.Help(); len(got) == 0 {
			t.Errorf("Help() is empty")
		}
	})
}

func TestRun(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
	}{
		{"empty", args{""}},
		{"simple", args{"simple"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := highlighter.Run(tt.args.input, nil); len(got) > 100 {
				t.Errorf("Run() = %v", got)
			}
		})
	}
}
