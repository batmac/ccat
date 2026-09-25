package mutators_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

// ~7 blocks at the minimum block size, with a short last block
var pbzip3Input = strings.Repeat("Lorem ipsum dolor sit amet, consectetur adipiscing elit. ", 8000)

func Test_pbzip3MatchesBzip3(t *testing.T) {
	for _, input := range []string{"", "x", "hello world", pbzip3Input} {
		want := mutators.Run("bzip3:66560", input)
		for _, w := range []int{0, 1, 2, 3, 8} {
			name := fmt.Sprintf("pbzip3:%d:66560", w)
			t.Run(fmt.Sprintf("%s/len%d", name, len(input)), func(t *testing.T) {
				if got := mutators.Run(name, input); got != want {
					t.Errorf("%s output differs from bzip3 (len %d vs %d)", name, len(got), len(want))
				}
			})
		}
	}
}

func Test_pbzip3DefaultBlockSize(t *testing.T) {
	if mutators.Run("pbzip3", "hello") != mutators.Run("bzip3", "hello") {
		t.Error("pbzip3 and bzip3 differ with default settings")
	}
	if mutators.Run("pbzip3:2", "hello") != mutators.Run("bzip3", "hello") {
		t.Error("pbzip3:2 and bzip3 differ with default block size")
	}
}

func Test_punbzip3(t *testing.T) {
	for _, input := range []string{"", "x", pbzip3Input} {
		compressed := mutators.Run("bzip3:66560", input)
		for _, m := range []string{"punbzip3", "punbzip3:1", "punbzip3:3", "unpbzip3:8"} {
			t.Run(fmt.Sprintf("%s/len%d", m, len(input)), func(t *testing.T) {
				if got := mutators.Run(m, compressed); got != input {
					t.Errorf("%s: len = %d, want %d", m, len(got), len(input))
				}
			})
		}
	}
}

func Test_pbzip3InvalidConfig(t *testing.T) {
	for _, name := range []string{"pbzip3:-1", "pbzip3:x", "pbzip3:2:1k", "pbzip3:1:1MB:3", "punbzip3:-1", "punbzip3:1:2"} {
		if _, err := mutators.New(name); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
}
