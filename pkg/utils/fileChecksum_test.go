package utils_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/utils"
)

func TestFileChecksum(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "empty_donotpanic",
			args:    args{path: ""},
			want:    "",
			wantErr: true,
		},
		{
			name:    "not found",
			args:    args{path: "notfound"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty file",
			args:    args{path: "testdata/empty"},
			want:    "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			wantErr: false,
		},
		{
			name:    "file",
			args:    args{path: "testdata/file"},
			want:    "e3ab5deea1aae346c0ae35d2aaf75b4a3557b88b40398b3f9d63ba15cc08fb1b",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.FileChecksum(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FileChecksum() = %v, want %v", got, tt.want)
			}
		})
	}
}
