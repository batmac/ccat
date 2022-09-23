package mutators_test

import (
	"strings"
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
	"github.com/batmac/ccat/pkg/stringutils"
)

var testsYJ = []struct {
	name, decoded, encoded string
}{
	{"null", "null", `null`},
	{"hi", `hello: hi`, `{"hello":"hi"}`},
	{"2xhi", "hello: hi\nhello2: hi", `{"hello":"hi","hello2":"hi"}`},
	{"parent", "first:\n  second: hi", `{"first":{"second":"hi"}}`},
	{"parent_", "first:\n  second:\n  - third: hi", `{"first":{"second":[{"third":"hi"}]}}`},
	{"list", "number:\n- 1\n- 2\n- 3", `{"number":[1,2,3]}`},
}

func Test_simpleY2J(t *testing.T) {
	f := "y2j"
	for _, tt := range testsYJ {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); stringutils.DeleteSpaces(got) != tt.encoded {
				t.Errorf("%s = %v, want %v", f, stringutils.DeleteSpaces(got), tt.encoded)
			}
		})
	}
}

func Test_simpleJ2Y(t *testing.T) {
	f := "j2y"
	for _, tt := range testsYJ {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.encoded); strings.TrimSuffix(got, "\n") != strings.TrimSuffix(tt.decoded, "\n") {
				t.Errorf("%s = %v, want %v", f, strings.TrimSuffix(got, "\n"), strings.TrimSuffix(tt.decoded, "\n"))
			}
		})
	}
}
