package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"mvdan.cc/gofumpt/format"
)

var (
	path   = "../../../LICENSE"
	target = "../generated_licence.go"
)

func main() {
	var b bytes.Buffer
	data, err := os.ReadFile(path)
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

	dataSize := len(data)
	gzDataSize := b.Len()

	gzData := printBuffer(&b)
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
	src, err := format.Source(b.Bytes(), format.Options{})
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(target, src, 0o600)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("size:%v bytes, compressed:%v bytes\n", dataSize, gzDataSize)
}

func printBuffer(data *bytes.Buffer) string {
	var count int
	var b strings.Builder
	b.Grow(data.Len()*6 + 9)
	b.WriteString("[]byte{\n")
	for _, byte := range data.Bytes() {
		fmt.Fprintf(&b, "0x%02x, ", byte)
		count++
		if count%12 == 0 {
			b.WriteString("\n")
		}
	}
	b.WriteString("}")
	return b.String()
}
