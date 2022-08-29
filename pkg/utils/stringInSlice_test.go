package utils_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/utils"
)

func TestStringInSlice(t *testing.T) {
	type args struct {
		a    string
		list []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty true", args{"", []string{"first", "notest", "cat", ""}}, true},
		{"empty false", args{"", []string{"first", "notest", "cat"}}, false},
		{"empty from empty -> false", args{"", []string{}}, false},
		{"test true", args{"test", []string{"first", "test", "notest", "cat", ""}}, true},
		{"test false", args{"test", []string{"first", "notest", "cat", ""}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.IsStringInSlice(tt.args.a, tt.args.list); got != tt.want {
				t.Errorf("StringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
