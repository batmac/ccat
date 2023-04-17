//go:build !fileonly
// +build !fileonly

package openers_test

import (
	"io"
	"testing"

	"github.com/batmac/ccat/pkg/openers"
)

func Test_echoOpener_Open(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		wantErr bool
		wanted  string
	}{
		{name: "ok", args: "echo://salut", wantErr: false, wanted: "salut"},
		{name: "ko", args: "echo://", wantErr: true, wanted: ""},
		{name: "ok double-quoted", args: "echo://\"salut yo\"", wantErr: false, wanted: "salut yo"},
		{name: "ok single-quoted", args: "echo://'salut again'", wantErr: false, wanted: "salut again"},
		{name: "ko double-quoted", args: "echo://\"\"", wantErr: true, wanted: ""},
		{name: "ko single-quoted", args: "echo://''", wantErr: true, wanted: ""},
		{name: "emoji", args: "echo://👍 🥳", wantErr: false, wanted: "👍 🥳"},
		{name: "with 'zero' byte in it", args: "echo://salut\x00and again", wantErr: false, wanted: "salut\x00and again"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := openers.Open(tt.args, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("echoOpener.Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil && !tt.wantErr {
				t.Errorf("echoOpener.Open() = got nil")
				return
			}
			if tt.wantErr {
				return
			}
			data, err := io.ReadAll(got)
			if err != nil {
				t.Errorf("echoOpener.Open() = %v", err)
				return
			}
			if string(data) != tt.wanted {
				t.Errorf("echoOpener.Open() = %#v, want %#v", string(data), tt.wanted)
				return
			}
		})
	}
}
