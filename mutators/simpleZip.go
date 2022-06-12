package mutators

import (
	"archive/zip"
	"bytes"
	"ccat/globalctx"
	"ccat/log"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
)

func init() {
	simpleRegister("unzip", unzip, withDescription("decompress the first file in a zip archive"))
	simpleRegister("zip", czip, withDescription("compress to zip data"), withExpectingBinary(true))
}

func unzip(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	dat, err := ioutil.ReadAll(in)
	if err != nil {
		log.Fatal(err)
	}
	// Open a zip archive for reading.
	r, err := zip.NewReader(bytes.NewReader(dat), int64(len(dat)))
	if err != nil {
		log.Fatal(err)
	}

	//return the first file
	for _, f := range r.File {
		log.Debugf("found file %s\n", f.Name)
		rc, err := f.Open()
		if err != nil {
			log.Debugln(err)
			continue
		}
		n, err := io.Copy(out, rc)
		if err != nil {
			log.Fatal(err)
		}
		rc.Close()
		return n, nil
	}
	return 0, fmt.Errorf("No extractable File found")
}

func czip(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	/* 	dat, err := ioutil.ReadAll(in)
	   	if err != nil {
	   		log.Fatal(err)
	   	} */
	z := zip.NewWriter(out)
	defer z.Close()

	f, err := z.Create(filename())
	if err != nil {
		log.Fatal(err)
	}
	return io.Copy(f, in)

}

func filename() string {
	name := filepath.Base(globalctx.Get("path").(string))
	return name
}
