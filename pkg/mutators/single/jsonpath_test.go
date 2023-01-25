package mutators_test

import (
	"strings"
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

var bookExample = `
	{
		"store":
		{
			"book":
			[
				{ "title": "Sayings of the Century", "author": "Nigel Rees", "price": 8.95 },
				{ "title": "Sword of Honour", "author": "Evelyn Waugh", "price": 12.99 }
			]
		}
	}
`

// see https://github.com/PaesslerAG/jsonpath/blob/master/jsonpath_test.go

func Test_simpleJsonpath(t *testing.T) {
	tests := []struct {
		name, json, jsonpath, result string
	}{
		{"hi", `{"hi":"hi"}`, ".hi", "\"hi\""},
		{"dot notation", bookExample, "$.store.book[0].title", "\"Sayings of the Century\""},
		{"current", `{"a":{"max":"3a", "3a":"aa"}, "1":{"a":"1a"}, "x":{"7":"bb"}}`, "$.a[@.max]", "\"aa\""},
		{"float equal", `{"a":1.23, "b":2}`, `$.a == 1.23`, "true"},
		{"ending star", `{"welcome":{"message":["Good Morning", "Hello World!"]}}`, `$.welcome.message[*]`, `["Good Morning","Hello World!"]`},
	}

	f := "jsonpath:"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f+tt.jsonpath, tt.json); strings.TrimSuffix(got, "\n") != tt.result {
				t.Errorf("%#v -> want %v, got %v", tt, tt.result, got)
			}
		})
	}
}
