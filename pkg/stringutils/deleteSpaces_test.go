package stringutils_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/stringutils"
)

func TestDeleteSpaces(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test", args{"test"}, "test"},
		{"test ", args{"test "}, "test"},
		{" test ", args{"test "}, "test"},
		{"  test  ", args{"  test  "}, "test"},
		{" te st ", args{"test "}, "test"},
		{"test\n", args{"test\n"}, "test"},
		{"\nte\nst\n", args{"\nte\nst\n"}, "test"},
		{"\ntest\n", args{"\ntest\n"}, "test"},
		{"te\t\n st", args{"te\t\n st"}, "test"},
		{"te  \u00a0st", args{"te  \u00a0st"}, "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringutils.DeleteSpaces(tt.args.str); got != tt.want {
				t.Errorf("DeleteSpaces() = %v, want %v", got, tt.want)
			}
		})
	}
}
