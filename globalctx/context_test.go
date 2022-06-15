package globalctx_test

import (
	"reflect"
	"testing"

	"github.com/batmac/ccat/globalctx"
)

var setDone bool
var testsGlobalCtx = []struct {
	name string
	k    string
	v    interface{}
}{
	{"bool", "bool", true},
	{"string", "string", "hi"},
	{"slice", "slice", []string{"a", "b", "c"}},
}

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
