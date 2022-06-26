package main

import (
	"testing"
)

func Test_update(t *testing.T) {
	t.Run("donotpanicplease", func(t *testing.T) {
		update("v0+dev", true)
		update("100", true)
	})
}

func Test_cleanVersion(t *testing.T) {
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
			want: "",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanVersion(tt.args.s); got != tt.want {
				t.Errorf("cleanVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
