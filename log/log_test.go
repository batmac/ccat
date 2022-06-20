//go:build !go1.15
// +build !go1.15

package log

import (
	"bytes"
	"log"
	"reflect"
	"strings"
	"testing"
)

func TestDefault(t *testing.T) {
	tests := []struct {
		name string
		want *Logger
	}{
		{"empty", &Logger{}},
		{"default", &Logger{log.Default()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Default(); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("Default() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetDebug(t *testing.T) {
	t.Run("bytes.Buffer", func(t *testing.T) {
		w := &bytes.Buffer{}
		SetDebug(w)
		if w != (Debug.Logger.Writer()) {
			t.Errorf("SetDebug() = %v, want %v", Debug.Logger, w)
		}
	})
}

func TestPp(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"null", args{nil}, "null"},
		{"zero", args{0}, "0"},
		{"empty", args{struct{}{}}, "{}"},
		{"some", args{"some"}, "\"some\""},
		{"slice", args{[]string{"0", "1"}}, "[\"0\",\"1\"]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(Pp(tt.args.data), "\n", ""), "\t", ""), " ", ""); got != tt.want {
				t.Errorf("Pp() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
