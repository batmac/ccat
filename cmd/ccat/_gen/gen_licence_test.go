package main

import (
	"bytes"
	"testing"
)

func Test_printBuffer(t *testing.T) {
	type args struct {
		data *bytes.Buffer
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// a few trivial test cases
		{
			name: "1",
			args: args{
				data: bytes.NewBuffer([]byte{}),
			},
			want: "[]byte{\n}",
		},
		{
			name: "2",
			args: args{
				data: bytes.NewBuffer([]byte{0x1, 0x2}),
			},
			want: "[]byte{\n0x01, 0x02, }",
		},
		{
			name: "3",
			args: args{
				data: bytes.NewBuffer([]byte{0x1, 0x2, 0x3}),
			},
			want: "[]byte{\n0x01, 0x02, 0x03, }",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := printBuffer(tt.args.data); got != tt.want {
				t.Errorf("printBuffer() = %v, want %v", got, tt.want)
			}
		})
	}
}
