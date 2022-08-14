package utils_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/utils"
)

func TestExpandPath(t *testing.T) {
	t.Setenv("TESTEXPANDPATH", "expandedpath")
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{""}, ""},
		{"simple", args{"/tmp"}, "/tmp"},
		{"simple", args{"/tmp/$NOTEXISTINGSTUFF"}, "/tmp/"},
		{"notStartingTilde", args{"/tmp/~/random"}, "/tmp/~/random"},
		{"TESTEXPANDPATH", args{"/tmp/$TESTEXPANDPATH"}, "/tmp/expandedpath"},
		{"TESTEXPANDPATHx2", args{"/tmp/$TESTEXPANDPATH/$TESTEXPANDPATH"}, "/tmp/expandedpath/expandedpath"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.ExpandPath(tt.args.path); got != tt.want {
				t.Errorf("ExpandPath() = %v, want %v", got, tt.want)
			}
		})
	}
	tests2 := []struct {
		name string
		args args
	}{
		{"PATH", args{"$PATH"}},
		{"tilde", args{"~"}},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.ExpandPath(tt.args.path); len(got) <= len(tt.args.path) {
				t.Errorf("ExpandPath() = %v, len = %v", got, len(got))
			}
		})
	}
}

func TestHomeDir(t *testing.T) {
	t.Run("HomeDir", func(t *testing.T) {
		if got := utils.HomeDir(); len(got) < 2 {
			t.Errorf("HomeDir() = %#v", got)
		}
	})
}
