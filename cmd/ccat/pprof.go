package main

import (
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/batmac/ccat/pkg/log"
	// _ "net/http/pprof"
)

var (
	cpuPprofFile *os.File
	memPprofFile *os.File
	err          error
)

func enablePprof() {
	log.Debugln("enabling pprof")

	cpuPprofFile, err = os.Create("cpu.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(cpuPprofFile); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}

	memPprofFile, err = os.Create("mem.pprof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
}

func endPprof() {
	log.Debugln("ending pprof")

	pprof.StopCPUProfile()
	_ = cpuPprofFile.Close()

	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(memPprofFile); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
	_ = memPprofFile.Close()
}
