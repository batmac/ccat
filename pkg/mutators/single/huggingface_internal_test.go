package mutators

import (
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

func Test_huggingface(t *testing.T) {
	// only test that we do not panic
	t.Setenv("CI", "CI")

	f := "huggingface"
	t.Run("donotpanicplease", func(t *testing.T) {
		if got := mutators.Run(f, "hi"); got != "fake" {
			t.Errorf("%s = %v, want %v", f, got, "fake")
		}
	})
}

func Test_getHuggingFaceToken(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		source  string
		wantErr bool
	}{
		{
			name:    "donotpanic",
			want:    "",
			source:  "",
			wantErr: true,
		},
	}
	t.Setenv("CI", "CI")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getHuggingFaceToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("getHuggingFaceToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getHuggingFaceToken() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.source {
				t.Errorf("getHuggingFaceToken() got1 = %v, want %v", got1, tt.source)
			}
		})
	}
}
