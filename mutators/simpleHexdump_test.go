package mutators_test

import (
	"testing"

	"github.com/batmac/ccat/mutators"
	"github.com/batmac/ccat/utils"
)

func Test_simpleHexDump(t *testing.T) {
	tests := []struct {
		name, input, want string
	}{
		{"hello", "hello", "00000000  68 65 6c 6c 6f  |hello|\n"},
		{"empty", "", ""},
		{"zero", "\x00", "00000000  00  |.|\n"},
		{
			"bytes",
			string([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}),
			"00000000  00 01 02 03 04 05 06 07  08 09 0a 0b 0c 0d 0e 0f  |................|  00000010  10 |.|\n",
		},
		{
			"empty", string([]byte{
				0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24,
				25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49,
				50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74,
				75, 76, 77, 78, 79,
			}),
			"00000000  00 01 02 03 04 05 06 07  08 09 0a 0b 0c 0d 0e 0f  |................|" +
				"00000010  10 11 12 13 14 15 16 17  18 19 1a 1b 1c 1d 1e 1f  |................|" +
				"00000020  20 21 22 23 24 25 26 27  28 29 2a 2b 2c 2d 2e 2f  | !\"#$%&'()*+,-./|" +
				"00000030  30 31 32 33 34 35 36 37  38 39 3a 3b 3c 3d 3e 3f  |0123456789:;<=>?|" +
				"00000040  40 41 42 43 44 45 46 47  48 49 4a 4b 4c 4d 4e 4f  |@ABCDEFGHIJKLMNO|",
		},
	}

	f := "hex"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.input); utils.DeleteSpaces(got) != utils.DeleteSpaces(tt.want) {
				t.Errorf("%s = %v, want %v", f, got, tt.want)
			}
		})
	}
}
