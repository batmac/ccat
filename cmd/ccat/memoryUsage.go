package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/batmac/ccat/pkg/utils"
)

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Fprintf(os.Stderr, "Alloc = %v ", utils.BytesForHumanString(m.Alloc))
	fmt.Fprintf(os.Stderr, "\tTotalAlloc = %v ", utils.BytesForHumanString(m.TotalAlloc))
	fmt.Fprintf(os.Stderr, "\tSys = %v ", utils.BytesForHumanString(m.Sys))
	fmt.Fprintf(os.Stderr, "\tNumGC = %v\n", m.NumGC)
}
