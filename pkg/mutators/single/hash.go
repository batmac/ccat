package mutators

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"

	"crypto/md5"  //#nosec
	"crypto/sha1" //#nosec

	"github.com/OneOfOne/xxhash"
	"github.com/zeebo/xxh3"
)

type hasher func() hash.Hash

var list = []struct {
	newHash     hasher
	name        string
	description string
}{
	{name: "sha256", description: "compute the sha256 checksum", newHash: sha256.New},
	// {"sha256", "compute the sha256 checksum", ssha256.New()},
	{name: "xxh64", description: "compute the xxhash64 checksum", newHash: xxh64New},
	{name: "xxh32", description: "compute the xxhash32 checksum", newHash: xxh32New},
	{name: "xxh3", description: "compute the xxh3 checksum", newHash: xxh3New},
	// useful but avoid them
	//#nosec
	{name: "md5", description: "compute the md5 checksum", newHash: md5.New},
	//#nosec
	{name: "sha1", description: "compute the sha1 checksum", newHash: sha1.New},
}

func init() {
	for _, h := range list {
		singlestRegister(h.name, wrap(h.newHash),
			withDescription(h.description),
			withCategory("checksum"))
	}
}

func wrap(f hasher) func(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	return func(w io.WriteCloser, r io.ReadCloser) (int64, error) {
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

func xxh3New() hash.Hash  { return xxh3.New() }
func xxh32New() hash.Hash { return xxhash.New32() }
func xxh64New() hash.Hash { return xxhash.New64() }
