package utils_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/utils"
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
			want: "1 B",
		},
		{
			name: "KB",
			args: args{
				b: 1000,
			},
			want: "1.00 kB",
		},
		{
			name: "MB",
			args: args{
				b: 1000 * 1000,
			},
			want: "1.00 MB",
		},
		{
			name: "GB",
			args: args{
				b: 1000 * 1000 * 1000,
			},
			want: "1.00 GB",
		},
		{
			name: "TB",
			args: args{
				b: 1000 * 1000 * 1000 * 1000,
			},
			want: "1.00 TB",
		},
		{
			name: "PB",
			args: args{
				b: 1000 * 1000 * 1000 * 1000 * 1000,
			},
			want: "1.00 PB",
		},
		{
			name: "EB",
			args: args{
				b: 1000 * 1000 * 1000 * 1000 * 1000 * 1000,
			},
			want: "1.00 EB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.BytesForHumanString(tt.args.b); got != tt.want {
				t.Errorf("BytesForHumanString() = %v, want %v", got, tt.want)
			}
		})
	}
}
