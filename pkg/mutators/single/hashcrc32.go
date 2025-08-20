package mutators

import (
	"encoding/hex"
	"hash/crc32"
	"io"
)

func init() {
	singleRegister("crc32", crc32Hash,
		withDescription("compute the crc32 checksum"),
		withCategory("checksum"))
}

func crc32Hash(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	h := crc32.NewIEEE()
	n, err := io.Copy(h, r)
	if err != nil {
		return 0, err
	}
	_, err = io.WriteString(w, hex.EncodeToString(h.Sum(nil))+"\n")
	if err != nil {
		return 0, err
	}
	return n, nil
}