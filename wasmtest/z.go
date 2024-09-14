package main

import (
	"bytes"
	"context"
	"fmt"
	// "io"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"
	// "github.com/tetratelabs/wazero/api"
)

func main() {
	ctx := context.Background()

	// Create a new WebAssembly runtime
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)

	// Instantiate WASI in the runtime
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	// Load your WASM module
	wasmBytes, err := os.ReadFile("main.wasm")
	if err != nil {
		panic(fmt.Sprintf("Failed to read WASM module: %v", err))
	}

	// Create an in-memory file system for stdin/stdout using bytes.Buffer
	stdin := new(bytes.Buffer)
	stdout := new(bytes.Buffer)

	// Write to stdin buffer (this simulates user input)
	stdin.WriteString("Hello, Wazero!\n")

	// Define the configuration with the in-memory file system
	config := wazero.NewModuleConfig().
		WithStdin(stdin).     // Set the in-memory stdin
		WithStdout(stdout).   // Set the in-memory stdout
		WithStderr(os.Stderr) // Set stderr to the standard error of the host process

	// Instantiate the WASM module with the configured input/output
	if _, err := r.InstantiateWithConfig(ctx, wasmBytes, config); err != nil {
		if exitErr, ok := err.(*sys.ExitError); ok && exitErr.ExitCode() != 0 {
			fmt.Fprintf(os.Stderr, "exit_code: %d\n", exitErr.ExitCode())
		} else if !ok {
			panic(err)
		}
	}

	// Run the "_start" function, which is the standard entry point for WASI modules
	// _, err = instance.ExportedFunction("_start").Call(ctx)
	// if err != nil {
	// panic(fmt.Sprintf("Failed to run WASM module: %v", err))
	// }

	// Read and print the output from the WASM module
	fmt.Printf("Output from WASM: %s", stdout.String())
}
