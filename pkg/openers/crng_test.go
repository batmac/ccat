//go:build !fileonly
// +build !fileonly

package openers //nolint

import (
	"io"
	"testing"
)

func Test_crngOpener_Open(t *testing.T) {
	tests := []struct {
		name         string
		args         string
		wantErr      bool
		wantedLength int
	}{
		{name: "ko", args: "crng://plop", wantErr: true},
		{name: "ko", args: "crng://-1", wantErr: true},
		{name: "ok", args: "crng://10", wantErr: false, wantedLength: 10},
		{name: "ok", args: "crng://1k", wantErr: false, wantedLength: 1000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := crngOpener{}
			got, err := f.Open(tt.args, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("crngOpener.Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil && !tt.wantErr {
				t.Errorf("crngOpener.Open() = got nil")
				return
			}
			if tt.wantErr {
				return
			}
			data, err := io.ReadAll(got)
			if err != nil {
				t.Errorf("crngOpener.Open() = %v", err)
				return
			}
			if len(data) != tt.wantedLength {
				t.Errorf("crngOpener.Open() = %v, want %v", len(data), tt.wantedLength)
				return
			}

			checkRandomness(t, data)
		})
	}
}

func checkRandomness(t *testing.T, data []byte) {
	t.Helper()
	acc := uint8(0)
	for i := 0; i < len(data); i++ {
		acc |= data[i]
	}
	if acc != 0xFF {
		t.Errorf("crngOpener.Open() = %v, want %v", acc, 255)
		return
	}
}

func Test_crngOpener_Evaluate(t *testing.T) {
	tests := []struct {
		name string
		args string
		want float32
	}{
		{name: "ko", args: "false://", want: 0},
		{name: "ok", args: "crng://", want: 0.9},
		{name: "ok with good limit", args: "crng://10", want: 0.9},
		{name: "ok with good limit 2", args: "crng://10k", want: 0.9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := crngOpener{}
			if got := f.Evaluate(tt.args); got != tt.want {
				t.Errorf("crngOpener.Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
