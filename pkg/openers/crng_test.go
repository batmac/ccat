//go:build !fileonly
// +build !fileonly

package openers_test

import (
	"io"
	"testing"

	"github.com/batmac/ccat/pkg/openers"
	"github.com/batmac/ccat/pkg/utils"
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
		{name: "ok", args: "crng://100", wantErr: false, wantedLength: 100},
		{name: "ok", args: "crng://1k", wantErr: false, wantedLength: 1000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := openers.Open(tt.args, false)
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

			utils.CheckBytesRandomness(t, data)
		})
	}
}
