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
	name, description string
	newHash           hasher
}{
	{"sha256", "compute the sha256 checksum", sha256.New},
	// {"sha256", "compute the sha256 checksum", ssha256.New()},
	{"xxh64", "compute the xxhash64 checksum", xxh64New},
	{"xxh32", "compute the xxhash32 checksum", xxh32New},
	{"xxh3", "compute the xxh3 checksum", xxh3New},
	// useful but avoid them
	//#nosec
	{"md5", "compute the md5 checksum", md5.New},
	//#nosec
	{"sha1", "compute the sha1 checksum", sha1.New},
}

func init() {
	for _, h := range list {
		simpleRegister(h.name, wrap(h.newHash),
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
