package mutators_test

import (
	"os"
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
	_ "github.com/batmac/ccat/pkg/mutators/single"
)

func Test_unlzfse(t *testing.T) {
	// .lzfse files can be compressed with lzfse or lzvn (for small files)
	// we test both
	paths := []string{
		"testdata/compression/lzfse.lzfse",
		"testdata/compression/lzvn.lzfse",
	}
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("Error reading file %s: %v", path, err)
		}
		// sanity check
		if len(content) < 100 {
			t.Errorf("Too small file %s: %v", path, len(content))
		}

		t.Run("unlzfse", func(t *testing.T) {
			uncompressed := mutators.Run("unlzfse", string(content))
			// sanity check
			if len(uncompressed) < len(content) {
				t.Errorf("unlzfse: len = %v (origin is %v)",
					len(uncompressed), len(content))
			}
		})
	}
}
