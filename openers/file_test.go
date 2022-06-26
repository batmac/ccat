package openers

import (
	"io"
	"reflect"
	"testing"
)

func Test_fileOpener_Name(t *testing.T) {
	type fields struct {
		name        string
		description string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fileOpener{
				name:        tt.fields.name,
				description: tt.fields.description,
			}
			if got := f.Name(); got != tt.want {
				t.Errorf("fileOpener.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileOpener_Description(t *testing.T) {
	type fields struct {
		name        string
		description string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fileOpener{
				name:        tt.fields.name,
				description: tt.fields.description,
			}
			if got := f.Description(); got != tt.want {
				t.Errorf("fileOpener.Description() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileOpener_Open(t *testing.T) {
	type fields struct {
		name        string
		description string
	}
	type args struct {
		s    string
		lock bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    io.ReadCloser
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fileOpener{
				name:        tt.fields.name,
				description: tt.fields.description,
			}
			got, err := f.Open(tt.args.s, tt.args.lock)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileOpener.Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fileOpener.Open() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileOpener_Evaluate(t *testing.T) {
	type fields struct {
		name        string
		description string
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fileOpener{
				name:        tt.fields.name,
				description: tt.fields.description,
			}
			if got := f.Evaluate(tt.args.s); got != tt.want {
				t.Errorf("fileOpener.Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parsePath(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parsePath(tt.args.s); got != tt.want {
				t.Errorf("parsePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
