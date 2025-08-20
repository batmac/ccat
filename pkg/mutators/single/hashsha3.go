package mutators

import (
	"crypto/sha3"
	"encoding/hex"
	"hash"
	"io"
)

type sha3hasher func() hash.Hash

var sha3list = []struct {
	newHash     sha3hasher
	name        string
	description string
}{
	{name: "sha3-224", description: "compute the sha3-224 checksum", newHash: func() hash.Hash { return sha3.New224() }},
	{name: "sha3-256", description: "compute the sha3-256 checksum", newHash: func() hash.Hash { return sha3.New256() }},
	{name: "sha3-384", description: "compute the sha3-384 checksum", newHash: func() hash.Hash { return sha3.New384() }},
	{name: "sha3-512", description: "compute the sha3-512 checksum", newHash: func() hash.Hash { return sha3.New512() }},
}

func init() {
	for _, h := range sha3list {
		singleRegister(h.name, wrapSha3(h.newHash),
			withDescription(h.description),
			withCategory("checksum"))
	}
}

func wrapSha3(f sha3hasher) func(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	return func(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
		h := f()
		h.Reset()
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
}