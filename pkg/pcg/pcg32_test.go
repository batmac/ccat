package pcg_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/pcg"
)

func TestNewPCG32(t *testing.T) {
	seedState := uint64(42)
	seedSeq := uint64(54)
	rng := pcg.NewPCG32(seedState, seedSeq)
	if rng == nil {
		t.Fatalf("NewPCG32 returned nil")
	}
	// Call Next() once to ensure initialization logic completes without panic
	_ = rng.Next()
}

func TestPCG32Output(t *testing.T) {
	seedState := uint64(42)
	seedSeq := uint64(54)
	rng := pcg.NewPCG32(seedState, seedSeq)

	// expectedValues := []uint32{... , ... , ...}
	// for i, expected := range expectedValues {
	//  	got := rng.Next()
	//  	if got != expected {
	//  		t.Errorf("Sequence mismatch at index %d: got %d, want %d", i, got, expected)
	//  	}
	// }

	// Test Case 2: Check for non-repeating values in short sequence
	val1 := rng.Next()
	val2 := rng.Next()
	val3 := rng.Next()

	if val1 == val2 || val2 == val3 {
		t.Errorf("Expected different values, got sequence: %d, %d, %d", val1, val2, val3)
	}

	// Test Case 3: Ensure values are generated (basic check)
	count := 100
	generated := make(map[uint32]struct{})
	var firstValue uint32
	allSame := true
	for i := 0; i < count; i++ {
		v := rng.Next()
		generated[v] = struct{}{}
		if i == 0 {
			firstValue = v
		} else if v != firstValue {
			allSame = false
		}
	}

	if len(generated) == 0 {
		t.Errorf("Generator produced no values.")
	}
	if allSame && count > 1 {
		t.Errorf("Generator produced the same value %d times: %d", count, firstValue)
	}
	if len(generated) < count/2 && count > 10 { // Expect at least some variety
		t.Errorf("Generated values seem non-random, only %d unique values in %d draws", len(generated), count)
	}

	t.Logf("Generated %d unique values in %d draws.", len(generated), count)
}

func BenchmarkPCG32Next(b *testing.B) {
	rng := pcg.NewPCG32(42, 54)
	for i := 0; i < b.N; i++ {
		_ = rng.Next()
	}
}
