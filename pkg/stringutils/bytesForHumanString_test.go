package stringutils_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/stringutils"
)

func TestBytesForHumanString(t *testing.T) {
	type args struct {
		b uint64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "B",
			args: args{
				b: 1,
			},
			want: "1B",
		},
		{
			name: "KB",
			args: args{
				b: 1000,
			},
			want: "1kB",
		},
		{
			name: "MB",
			args: args{
				b: 1000 * 1000,
			},
			want: "1MB",
		},
		{
			name: "GB",
			args: args{
				b: 1000 * 1000 * 1000,
			},
			want: "1GB",
		},
		{
			name: "TB",
			args: args{
				b: 1000 * 1000 * 1000 * 1000,
			},
			want: "1TB",
		},
		{
			name: "PB",
			args: args{
				b: 1000 * 1000 * 1000 * 1000 * 1000,
			},
			want: "1PB",
		},
		{
			name: "EB",
			args: args{
				b: 1000 * 1000 * 1000 * 1000 * 1000 * 1000,
			},
			want: "1EB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringutils.BytesForHumanString(tt.args.b); got != tt.want {
				t.Errorf("BytesForHumanString() = %v, want %v", got, tt.want)
			}
		})
	}
}
