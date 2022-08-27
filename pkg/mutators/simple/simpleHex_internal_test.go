package mutators

import (
	"testing"
)

func Test_isBase16Char(t *testing.T) {
	type args struct {
		c byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "0", args: args{'0'}, want: true},
		{name: "9", args: args{'9'}, want: true},
		{name: "A", args: args{'A'}, want: true},
		{name: "F", args: args{'F'}, want: true},
		{name: "G", args: args{'G'}, want: false},
		{name: "Z", args: args{'Z'}, want: false},
		{name: "a", args: args{'a'}, want: true},
		{name: "f", args: args{'f'}, want: true},
		{name: "g", args: args{'g'}, want: false},
		{name: "z", args: args{'z'}, want: false},
		{name: " ", args: args{' '}, want: false},
		{name: "nl", args: args{'\n'}, want: false},
		{name: "tab", args: args{'\t'}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isBase16Char(tt.args.c); got != tt.want {
				t.Errorf("isBase16Char(%s) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
