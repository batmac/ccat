package log_test

import (
	"bytes"
	"strings"
	"testing"

	// stdlog "log"

	"github.com/batmac/ccat/pkg/log"
)

/* func TestDefault(t *testing.T) {
	tests := []struct {
		name string
		want *log.Logger
	}{
		{"empty", &log.Logger{}},
		{"default", &log.Logger{stdlog.Default()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := log.Default(); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("Default() = %v, want %v", got, tt.want)
			}
		})
	}
} */

func TestSetDebug(t *testing.T) {
	t.Run("bytes.Buffer", func(t *testing.T) {
		w := &bytes.Buffer{}
		log.SetDebug(w)
		if w != (log.Debug.Logger.Writer()) {
			t.Errorf("SetDebug() = %v, want %v", log.Debug.Logger, w)
		}
	})
}

func TestPp(t *testing.T) {
	type args struct {
		data any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"null", args{nil}, "null"},
		{"zero", args{0}, "0"},
		{"empty", args{struct{}{}}, "{}"},
		{"some", args{"some"}, "\"some\""},
		{"slice", args{[]string{"0", "1"}}, "[\"0\",\"1\"]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(log.Pp(tt.args.data), "\n", ""), "\t", ""), " ", ""); got != tt.want {
				t.Errorf("Pp() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestDontPanic(t *testing.T) {
	_ = log.Default()
	t.Log("Start of test")
	log.Debugf("hi")
	log.Debugln("hi")
	log.Print("hi")
	log.Printf("hi")
	log.Println("hi")
	log.SetContinueOnFatal()
	log.Fatal("hi")
}
