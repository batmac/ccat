package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func TestCompressionGo() error {
	stepPrintln("Testing compression...")

	failure := false

	dir := filepath.FromSlash("testdata/compression")
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		filePath := filepath.Join(dir, file.Name())
		expectedCksum, err := fileChecksum(filePath)
		if err != nil {
			return err
		}
		if mg.Debug() {
			fmt.Printf("%v %v\n", expectedCksum, file.Name())
		}
		for _, alg := range compressionAlgs() {
			opts := []string{"-m", alg + ",un" + alg + ",sha256", filePath}
			cksum, err := sh.Output("./"+binaryName, opts...)
			if err != nil {
				return err
			}
			if cksum != expectedCksum {
				failure = true
				if mg.Debug() {
					fmt.Printf("%v (%v) failed ! (%v != %v)\n",
						file.Name(), alg, cksum, expectedCksum)
				}
			}
		}
	}
	if failure {
		return errors.New("some checksum(s) don't match")
	}
	stepOKPrintln("Testing compression OK")
	return nil
}

func fileChecksum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func compressionAlgs() []string {
	// this should be automated
	algs := "bzip2 gzip lz4 lzma2 lzma s2 snap xz zip zlib zstd"
	return strings.Split(algs, " ")
}
