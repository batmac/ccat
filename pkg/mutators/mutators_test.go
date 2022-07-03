package mutators_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
	_ "github.com/batmac/ccat/pkg/mutators/simple"
)

// yo dawg
func Test_Run(t *testing.T) {
	type args struct {
		f     string
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "first",
			args: args{
				f:     "dummy",
				input: "hello",
			},
			want: "hello",
		},
		{
			name: "empty",
			args: args{
				f:     "dummy",
				input: "",
			},
			want: "",
		},
		{
			name: "zero",
			args: args{
				f:     "dummy",
				input: "\x00",
			},
			want: "\x00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(tt.args.f, tt.args.input); got != tt.want {
				t.Errorf("mutators.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
