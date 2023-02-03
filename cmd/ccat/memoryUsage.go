package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/docker/go-units"
)

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Fprintf(os.Stderr, "Alloc = %v ", units.HumanSize(float64(m.Alloc)))
	fmt.Fprintf(os.Stderr, "\tTotalAlloc = %v ", units.HumanSize(float64(m.TotalAlloc)))
	fmt.Fprintf(os.Stderr, "\tSys = %v ", units.HumanSize(float64(m.Sys)))
	fmt.Fprintf(os.Stderr, "\tNumGC = %v\n", m.NumGC)
}
