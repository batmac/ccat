package globalctx_test

import (
	"reflect"
	"testing"

	"github.com/batmac/ccat/pkg/globalctx"
)

var (
	setDone        bool
	testsGlobalCtx = []struct {
		name string
		k    string
		v    interface{}
	}{
		{"bool", "bool", true},
		{"string", "string", "hi"},
		{"slice", "slice", []string{"a", "b", "c"}},
	}
)

func TestSet(t *testing.T) {
	for _, tt := range testsGlobalCtx {
		t.Run(tt.name, func(t *testing.T) {
			globalctx.Set(tt.k, tt.v)
		})
	}
	setDone = true
}

func TestGet(t *testing.T) {
	if !setDone {
		TestSet(t)
	}
	for _, tt := range testsGlobalCtx {
		t.Run(tt.name, func(t *testing.T) {
			got := globalctx.Get(tt.k)
			if !reflect.DeepEqual(got, tt.v) {
				t.Errorf("Get() = %v, want %v", got, tt.v)
			}
		})
	}
}

func TestReset(t *testing.T) {
	if !setDone {
		TestSet(t)
	}
	globalctx.Reset()
	for _, tt := range testsGlobalCtx {
		t.Run(tt.name, func(t *testing.T) {
			got := globalctx.Get(tt.k)
			if reflect.DeepEqual(got, tt.v) {
				t.Errorf("Get() = %v, want nil", got)
			}
		})
	}
}

func TestIsErrored(t *testing.T) {
	if !setDone {
		TestSet(t)
	}
	globalctx.SetErrored()
	if !globalctx.IsErrored() {
		t.Errorf("IsErrored = false")
	}
	globalctx.Reset()
	if !globalctx.IsErrored() {
		t.Errorf("after Reset, IsErrored = false")
	}
}

func TestGetBool(t *testing.T) {
	globalctx.Set("test_true", true)
	globalctx.Set("test_false", false)
	type args struct {
		k string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"true", args{"test_true"}, true},
		{"false", args{"test_false"}, false},
		{"not existing", args{"not existing"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := globalctx.GetBool(tt.args.k); got != tt.want {
				t.Errorf("GetBool() = %v, want %v", got, tt.want)
			}
		})
	}
}
