package mutators

import (
	"reflect"
	"testing"
)

func Test_stdConfigHumanSizeAsInt64(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    int64
		wantErr bool
	}{
		{
			name:    "empty string",
			args:    []string{""},
			wantErr: true,
		},
		{
			name:    "empty",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "invalid",
			args:    []string{"abc"},
			wantErr: true,
		},

		{
			name:    "simple",
			args:    []string{"10"},
			want:    10,
			wantErr: false,
		},
		{
			name:    "with unit",
			args:    []string{"10k"},
			want:    10000,
			wantErr: false,
		},
		{
			name:    "two",
			args:    []string{"foo", "bar"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stdConfigHumanSizeAsInt64(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("stdConfigHumanSizeAsInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stdConfigHumanSizeAsInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stdConfigString(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{
			name:    "empty",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "simple",
			args:    []string{"abc"},
			want:    "abc",
			wantErr: false,
		},

		{
			name:    "number",
			args:    []string{"10"},
			want:    "10",
			wantErr: false,
		},
		{
			name:    "two",
			args:    []string{"foo", "bar"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stdConfigString(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("stdConfigString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stdConfigString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stdConfigStringsAtLeast(t *testing.T) {
	tests := []struct {
		name    string
		atLeast int
		upTo    int
		args    []string
		want    []string
		wantErr bool
	}{
		{
			name:    "empty",
			atLeast: 1,
			upTo:    1,
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "simple",
			atLeast: 1,
			upTo:    1,
			args:    []string{"abc"},
			want:    []string{"abc"},
			wantErr: false,
		},

		{
			name:    "number",
			atLeast: 1,
			upTo:    1,
			args:    []string{"10"},
			want:    []string{"10"},
			wantErr: false,
		},
		{
			name:    "less than two",
			atLeast: 2,
			upTo:    2,
			args:    []string{"foo"},
			wantErr: true,
		},
		{
			name:    "two",
			atLeast: 2,
			upTo:    2,
			args:    []string{"foo", "bar"},
			want:    []string{"foo", "bar"},
			wantErr: false,
		},
		{
			name:    "more than two",
			atLeast: 2,
			upTo:    4,
			args:    []string{"foo", "bar", "baz", "qux"},
			want:    []string{"foo", "bar", "baz", "qux"},
			wantErr: false,
		},
		{
			name:    "too many",
			atLeast: 2,
			upTo:    3,
			args:    []string{"foo", "bar", "baz", "qux"},
			want:    []string{"foo", "bar", "baz", "qux"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := stdConfigStrings(tt.atLeast, tt.upTo)
			res, err := fn(tt.args)
			var got []string
			var ok bool

			if (err != nil) != tt.wantErr {
				t.Errorf("stdConfigStringsAtLeast() error = %#v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got, ok = res.([]string); !ok {
				t.Errorf("stdConfigStringsAtLeast() = %#v, want []string", res)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stdConfigStringsAtLeast() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
