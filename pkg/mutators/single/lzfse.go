package mutators

import (
	"bytes"
	"io"

	"github.com/batmac/ccat/pkg/log"

	lzfse_go "github.com/aixiansheng/lzfse"
	// "github.com/blacktop/lzfse-cgo"
)

func init() {
	// singleRegister("unlzfse-c", unlzfse, withDescription("decompress lzfse data"))
	singleRegister("unlzfse", unlzfseGo, withDescription("decompress lzfse data"),
		withCategory("decompress"),
	)
	// singleRegister("lzfse", clzfse, withDescription("compress lzfse data"))
}

/* func unlzfse(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	dat, err := io.ReadAll(in)
	if err != nil {
		log.Fatal("failed to read compressed file: ", err)
	}
	decompressed := lzfse.DecodeBuffer(dat)
	d := bytes.NewReader(decompressed)
	return io.Copy(out, d)
} */

func unlzfseGo(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	dat, err := io.ReadAll(in) // NOT streamable
	if err != nil {
		log.Fatal("failed to read compressed file: ", err)
	}

	d := lzfse_go.NewReader(bytes.NewReader(dat))
	return io.Copy(out, d)
}
