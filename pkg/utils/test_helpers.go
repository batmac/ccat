package utils

import "testing"

func CheckBytesRandomness(t *testing.T, data []byte) {
	t.Helper()
	acc := uint8(0)
	for i := 0; i < len(data); i++ {
		acc |= data[i]
	}
	if acc != 0xFF {
		t.Errorf("CheckBytesRandomness(): %v, want %v", acc, 255)
		return
	}
}
