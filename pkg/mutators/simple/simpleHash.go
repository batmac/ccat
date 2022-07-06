package mutators

import (
	"encoding/hex"
	"hash"
	"io"

	ootb "crypto/sha256"

	"github.com/minio/sha256-simd"

	"github.com/cespare/xxhash/v2"
)

func init() {
	simpleRegister("sha256std", base, withDescription("compute the sha256 (from stdlib) checksum"),
		withCategory("checksum"))
	simpleRegister("sha256", cksum256, withDescription("compute the sha256 checksum"),
		withCategory("checksum"))
	simpleRegister("xxhash", cksumxxh64, withDescription("compute the xxhash (64bits) checksum"),
		withCategory("checksum"))
}

func base(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	return cksum(ootb.New(), w, r)
}

func cksum256(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	return cksum(sha256.New(), w, r)
}

func cksumxxh64(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	return cksum(xxhash.New(), w, r)

}

func cksum(h hash.Hash, w io.WriteCloser, r io.ReadCloser) (int64, error) {
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
