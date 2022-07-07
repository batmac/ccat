package mutators

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"

	"crypto/md5"  //#nosec
	"crypto/sha1" //#nosec

	"github.com/cespare/xxhash/v2"
	//  ssha256 "github.com/minio/sha256-simd"
	"github.com/zeebo/xxh3"
)

var list = []struct {
	name, description string
	hash              hash.Hash
}{
	{"sha256", "compute the sha256 (stdlib) checksum", sha256.New()},
	// {"sha256", "compute the sha256 checksum", ssha256.New()},
	{"xxhash", "compute the xxhash (xxh64) checksum", xxhash.New()},
	{"xxh3", "compute the xxh3 checksum", xxh3.New()},
	// useful but avoid them
	//#nosec
	{"md5", "compute the md5 checksum", md5.New()},
	//#nosec
	{"sha1", "compute the sha1 checksum", sha1.New()},
}

func init() {
	for _, h := range list {
		simpleRegister(h.name, wrap(h.hash),
			withDescription(h.description),
			withCategory("checksum"))
	}
}

func wrap(h hash.Hash) func(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	return func(w io.WriteCloser, r io.ReadCloser) (int64, error) {
		h.Reset()
		n, err := io.Copy(h, r)
		if err != nil {
			return 0, err
		}
		_, err = w.Write([]byte(hex.EncodeToString(h.Sum(nil)) + "\n"))
		if err != nil {
			return 0, err
		}
		return n, nil
	}
}
