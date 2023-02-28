//go:build !fileonly
// +build !fileonly

package openers_test

import (
	"io"
	"testing"

	"github.com/batmac/ccat/pkg/openers"
	"github.com/batmac/ccat/pkg/utils"
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
		{name: "ok", args: "prng://", wantErr: false, wantedLength: 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := openers.Open(tt.args, false)
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

			utils.CheckBytesRandomness(t, data)
		})
	}
}
