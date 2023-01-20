//go:build !fileonly
// +build !fileonly

package openers //nolint

import (
	"io"
	"testing"
)

func Test_prngOpener_Open(t *testing.T) {
	tests := []struct {
		name         string
		args         string
		wantErr      bool
		wantedLength int
	}{
		{name: "ko1", args: "prng://plop", wantErr: true},
		{name: "ko2", args: "prng://1k", wantErr: true},
		{name: "ok", args: "prng://", wantErr: false, wantedLength: 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := prngOpener{}
			got, err := f.Open(tt.args, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("prngOpener.Open() %s: error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}

			if got == nil && !tt.wantErr {
				t.Errorf("prngOpener.Open() = got nil")
				return
			}
			if tt.wantErr {
				return
			}
			data := make([]byte, tt.wantedLength)
			n, err := io.ReadFull(got, data)
			if err != nil {
				t.Errorf("prngOpener.Open() = %v", err)
				return
			}
			if n != tt.wantedLength {
				t.Errorf("prngOpener.Open() = %v, want %v", len(data), tt.wantedLength)
				return
			}

			checkRandomness(t, data)
		})
	}
}

func Test_prngOpener_Evaluate(t *testing.T) {
	tests := []struct {
		name string
		args string
		want float32
	}{
		{name: "ko", args: "false://", want: 0},
		{name: "ok", args: "prng://", want: 0.9},
		{name: "ok with good limit", args: "prng://10", want: 0.9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := prngOpener{}
			if got := f.Evaluate(tt.args); got != tt.want {
				t.Errorf("prngOpener.Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
