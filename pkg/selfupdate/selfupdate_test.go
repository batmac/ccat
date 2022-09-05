package selfupdate_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/selfupdate"
)

func Test_Do(t *testing.T) {
	t.Run("donotpanicplease", func(t *testing.T) {
		selfupdate.Do("v0+dev", "", true)
		selfupdate.Do("100", "", true)
	})
}

func Test_CleanVersion(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				s: "",
			},
			want: "0",
		},
		{
			name: "1",
			args: args{
				s: "1",
			},
			want: "1",
		},
		{
			name: "v1.9",
			args: args{
				s: "v1.9",
			},
			want: "1.9",
		},
		{
			name: "dev",
			args: args{
				s: ">0.9.8+dev",
			},
			want: "0.9.8",
		},
		{
			name: "stuff",
			args: args{
				s: "hpo1uoyru.sd0es_ 😀",
			},
			want: "1.0",
		},
		{
			name: "patch",
			args: args{
				s: "0.9.9-11-a23",
			},
			want: "0.9.9",
		},
		{
			name: "vpatch",
			args: args{
				s: "v0.9.9-11-a23",
			},
			want: "0.9.9",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := selfupdate.CleanVersion(tt.args.s); got != tt.want {
				t.Errorf("CleanVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
