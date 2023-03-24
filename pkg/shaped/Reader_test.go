package shaped_test

import (
	"io"
	"strings"
	"testing"

	"github.com/batmac/ccat/pkg/shaped"
)

func TestReader_Read(t *testing.T) {
	type fields struct {
		r         io.Reader
		bandwidth int
	}

	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "salut",
			fields: fields{
				r:         strings.NewReader("salut"),
				bandwidth: 100,
			},
			want:    "salut",
			wantErr: false,
		},
		{
			name: "hi",
			fields: fields{
				r:         strings.NewReader("hi"),
				bandwidth: 1,
			},
			want:    "hi",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lr := shaped.NewReader(tt.fields.r, tt.fields.bandwidth)

			buf := strings.Builder{}
			_, err := io.Copy(&buf, lr)
			got := buf.String()
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Reader.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
