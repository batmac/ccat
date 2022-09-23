package mutators_test

import (
	"os"
	"testing"

	"github.com/batmac/ccat/pkg/globalctx"
	"github.com/batmac/ccat/pkg/mutators"
	_ "github.com/batmac/ccat/pkg/mutators/simple"
)

func TestCompressionAlgs(t *testing.T) {
	path := "testdata/compression/test.txt"
	globalctx.Set("path", path)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("Error reading file %s: %v", path, err)
	}
	algs := mutators.ListAvailableMutators("compress")
	// sanity check
	if len(algs) < 10 {
		t.Errorf("Too few compression algorithms: %v", algs)
	}
	for _, f := range algs {
		t.Run(f, func(t *testing.T) {
			compressed := mutators.Run(f, string(content))
			// sanity check
			if len(compressed) >= len(content) || len(compressed) <= len(content)/10 {
				t.Errorf("%s: len = %v (origin is %v)", f, len(compressed), len(content))
			}
			uncompressed := mutators.Run("un"+f, compressed)
			// sanity check
			if len(uncompressed) != len(content) {
				t.Errorf("%s: len = %v (origin is %v)", f, len(uncompressed), len(content))
			}
			if uncompressed != string(content) {
				t.Errorf("%s: content != uncompressed", f)
			}
		})
	}
}
