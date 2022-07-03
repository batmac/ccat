package lockable_test

import (
	"os"
	"testing"

	"github.com/batmac/ccat/pkg/lockable"
)

func TestFileOpen(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		path string
		lock bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"not existing", args{"not existing", false}, true},
		{"not existing", args{"not existing", true}, true},
		{"", args{exe, false}, false},
		{"", args{exe, true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := lockable.FileOpen(tt.args.path, tt.args.lock)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileOpen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
