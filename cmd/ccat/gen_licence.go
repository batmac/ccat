//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"text/template"

	"github.com/batmac/ccat/pkg/mutators"
	_ "github.com/batmac/ccat/pkg/mutators/simple"
)

var (
	path   = "../../LICENSE"
	target = "generated_licence.go"
)

func main() {
	var b bytes.Buffer
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	gz, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		log.Fatal(err)
	}
	_, err = gz.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	if err = gz.Flush(); err != nil {
		log.Fatal(err)
	}
	if err = gz.Close(); err != nil {
		log.Fatal(err)
	}

	/* 	if err = ioutil.WriteFile(path+".gz", b.Bytes(), 0644); err != nil {
		log.Fatal(err)
	} */

	gzData := fmt.Sprintf("%#v\n", b.Bytes())
	gzData = mutators.Run("wrap", gzData)
	b.Reset()
	err = template.Must(template.New("").Parse(`
	// Code generated automatically. DO NOT EDIT.
	package main
	import (
		"bytes"
		"compress/gzip"
		"io"
		"log"
	)
	func printLicense(w io.Writer) {
	    var data = {{ . }}
		zr, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			log.Fatal(err)
		}
		if _, err := io.CopyN(w, zr, 10*1024); err != io.EOF {
			log.Fatal(err)
		}
		if err := zr.Close(); err != nil {
			log.Fatal(err)
		}
	}
	`)).Execute(&b, gzData)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(target, b.Bytes(), 0o644)

	cmd := exec.Command("gofmt", "-s", "-w", target)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(err)
	}
	if len(stdoutStderr) > 0 {
		log.Printf("%s\n", stdoutStderr)
	}
}
