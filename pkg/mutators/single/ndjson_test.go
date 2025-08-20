package mutators_test

import (
	"strings"
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

var testsNDJSONIndent = []struct {
	name, input, expected string
}{
	{
		"simple",
		`{"name":"John","age":30}`,
		`{
  "name": "John",
  "age": 30
}`,
	},
	{
		"multiple_lines",
		`{"name":"John","age":30}
{"name":"Jane","age":25}`,
		`{
  "name": "John",
  "age": 30
}
{
  "name": "Jane",
  "age": 25
}`,
	},
	{
		"empty_line",
		`{"name":"John","age":30}

{"name":"Jane","age":25}`,
		`{
  "name": "John",
  "age": 30
}

{
  "name": "Jane",
  "age": 25
}`,
	},
	{
		"invalid_json",
		`{"name":"John","age":30}
invalid json line
{"name":"Jane","age":25}`,
		`{
  "name": "John",
  "age": 30
}
invalid json line
{
  "name": "Jane",
  "age": 25
}`,
	},
	{
		"nested_object",
		`{"user":{"name":"John","profile":{"age":30,"city":"NYC"}}}`,
		`{
  "user": {
    "name": "John",
    "profile": {
      "age": 30,
      "city": "NYC"
    }
  }
}`,
	},
	{
		"array",
		`{"items":[1,2,3]}`,
		`{
  "items": [
    1,
    2,
    3
  ]
}`,
	},
}

func Test_NDJSONIndent(t *testing.T) {
	f := "ndjsonindent"
	for _, tt := range testsNDJSONIndent {
		t.Run(tt.name, func(t *testing.T) {
			got := mutators.Run(f, tt.input)
			// Remove trailing newline for comparison
			got = strings.TrimSuffix(got, "\n")
			if got != tt.expected {
				t.Errorf("%s = %v, want %v", f, got, tt.expected)
			}
		})
	}
}