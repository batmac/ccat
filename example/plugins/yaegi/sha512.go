package plugins

import (
	"crypto/sha512"
	"encoding/hex"
	"io"
)

// ccat -m yaegi:example/plugins/sha512.go:plugins.Sha512 go.mod

func Sha512(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	// compute the sha512 of r
	h := sha512.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return 0, err
	}
	n, err := io.WriteString(w, hex.EncodeToString(h.Sum(nil))+"\n")
	return int64(n), err
}
