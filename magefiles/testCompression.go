package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"braces.dev/errtrace"
	"github.com/batmac/ccat/pkg/mutators"
	_ "github.com/batmac/ccat/pkg/mutators/single"
	"github.com/batmac/ccat/pkg/utils"
)

func TestCompressionGo() error {
	stepPrintln("Testing compression...")

	failure := false

	dir := filepath.FromSlash("testdata/compression")
	files, err := os.ReadDir(dir)
	if err != nil {
		return errtrace.Wrap(err)
	}
	for _, file := range files {
		filePath := filepath.Join(dir, file.Name())
		expectedCksum, err := utils.FileChecksum(filePath)
		if err != nil {
			return errtrace.Wrap(err)
		}
		if mg.Verbose() {
			fmt.Printf("%v %v\n", expectedCksum, file.Name())
		}
		for _, alg := range compressionAlgs() {
			opts := []string{"-m", alg + ",un" + alg + ",sha256", filePath}
			cksum, err := sh.Output("./"+binaryName, opts...)
			if err != nil {
				return errtrace.Wrap(err)
			}
			if mg.Debug() {
				fmt.Printf("%v %v\n", cksum, alg)
			}
			if cksum != expectedCksum {
				failure = true
				if mg.Verbose() {
					fmt.Printf("%v (%v) failed ! (%v != %v)\n",
						file.Name(), alg, cksum, expectedCksum)
				}
			}
		}
	}
	if failure {
		return errtrace.Wrap(errors.New("some checksum(s) don't match"))
	}
	stepOKPrintln("Testing compression OK")
	return nil
}

func compressionAlgs() []string {
	algs := mutators.ListAvailableMutators("compress")
	// sanity check
	if len(algs) < 10 {
		panic("not enough compression algorithms")
	}
	return algs
}
