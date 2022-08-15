package main

import (
	"io"
	"os"

	"github.com/batmac/ccat/pkg/globalctx"
	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/openers"
)

func processFileAsIs(w io.Writer, s string) {
	var from io.ReadCloser = os.Stdin
	var err error
	if s != "-" {
		from, err = openers.Open(s, false)
		if err != nil {
			globalctx.SetErrored()
			log.Println(err)
			return
		}
		defer from.Close()
	}
	_, err = io.Copy(w, from)
	if err != nil {
		globalctx.SetErrored()
		log.Println(err)
		return
	}
	log.Debugf("end %s\n", s)
}
