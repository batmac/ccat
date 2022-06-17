package utils

import "testing"

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
			if got := StringInSlice(tt.args.a, tt.args.list); got != tt.want {
				t.Errorf("StringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			if got := DeleteSpaces(tt.args.str); got != tt.want {
				t.Errorf("DeleteSpaces() = %v, want %v", got, tt.want)
			}
		})
	}
}
