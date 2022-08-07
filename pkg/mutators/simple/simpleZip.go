package mutators

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"path/filepath"

	"github.com/batmac/ccat/pkg/globalctx"
	"github.com/batmac/ccat/pkg/log"
)

func init() {
	simpleRegister("unzip", unzip, withDescription("decompress the first file in a zip archive"),
		withCategory("decompress"),
	)
	simpleRegister("zip", czip, withDescription("compress to zip data"),
		withCategory("compress"),
	)
}

func unzip(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	dat, err := io.ReadAll(in)
	if err != nil {
		log.Fatal(err)
	}
	// Open a zip archive for reading.
	r, err := zip.NewReader(bytes.NewReader(dat), int64(len(dat)))
	if err != nil {
		log.Fatal(err)
	}

	// return the first file
	for _, f := range r.File {
		log.Debugf("found file %s\n", f.Name)
		rc, err := f.Open()
		if err != nil {
			log.Debugln(err)
			continue
		}
		//#nosec
		n, err := io.Copy(out, rc)
		if err != nil {
			log.Fatal(err)
		}
		rc.Close()
		return n, nil
	}
	return 0, fmt.Errorf("no extractable File found")
}

func czip(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	/* 	dat, err := io.ReadAll(in)
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
