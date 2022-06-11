package mutators

import (
	"bytes"
	"ccat/log"
	"io"
	"io/ioutil"

	lzfse_go "github.com/aixiansheng/lzfse"
	//"github.com/blacktop/lzfse-cgo"
)

func init() {
	//simpleRegister("unlzfse-c", unlzfse, withDescription("decompress lzfse data"))
	simpleRegister("unlzfse", unlzfseGo, withDescription("decompress lzfse data"))
	//simpleRegister("lzfse", clzfse, withDescription("compress lzfse data"))
}

/* func unlzfse(out io.WriteCloser, in io.ReadCloser) (int64, error) {

dat, err := ioutil.ReadAll(in)

if err != nil {
	log.Fatal("failed to read compressed file: ", err)
}

decompressed := lzfse.DecodeBuffer(dat)

d := bytes.NewReader(decompressed)

return io.Copy(out, d) */

func unlzfseGo(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	dat, err := ioutil.ReadAll(in)
	if err != nil {
		log.Fatal("failed to read compressed file: ", err)
	}

	d := lzfse_go.NewReader(bytes.NewReader(dat))
	return io.Copy(out, d)
}

/* func clzfse(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	dat, err := ioutil.ReadAll(in)
	if err != nil {
		log.Fatal("failed to read file: ", err)
	}

	b := make([]byte, len(dat), len(dat))
	s := lzfse.EncodeLZVNBuffer(dat, b)
	log.Debugf("EncodeLZVNBuffer %v\n", s)
	return io.Copy(out, bytes.NewReader(b))

}
*/
