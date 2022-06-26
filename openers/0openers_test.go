package openers_test

import (
	"os"
	"testing"

	"github.com/batmac/ccat/openers"
)

func TestListOpenersWithDescription(t *testing.T) {
	t.Run("notempty", func(t *testing.T) {
		if got := openers.ListOpenersWithDescription(); len(got) == 0 {
			t.Errorf("ListOpenersWithDescription() is empty")
		}
	})
}

func TestOpen(t *testing.T) {
	exe, err := os.Executable() // we just need an existing file
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		s    string
		lock bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"not existing", args{"not existing", false}, true},
		{"not existing", args{"not existing", true}, true},
		{"not existing", args{"fake://" + "not existing", false}, true},

		{"exe", args{exe, false}, false},
		{"exe", args{exe, true}, false},
		{"exe", args{"file://" + exe, false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := openers.Open(tt.args.s, tt.args.lock)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
