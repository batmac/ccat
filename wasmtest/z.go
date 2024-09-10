package main

import (
	"bytes"
	"context"
	"fmt"
	// "io"
	"os"

	"github.com/tetratelabs/wazero"
	// "github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func main() {
	ctx := context.Background()

	// Create a new WebAssembly runtime
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)

	// Instantiate WASI in the runtime
	_, err := wasi_snapshot_preview1.Instantiate(ctx, r)
	if err != nil {
		panic(fmt.Sprintf("Failed to instantiate WASI: %v", err))
	}

	// Load your WASM module
	wasmBytes, err := os.ReadFile("main.wasm")
	if err != nil {
		panic(fmt.Sprintf("Failed to read WASM module: %v", err))
	}

	// Compile the WASM module
	compiled, err := r.CompileModule(ctx, wasmBytes)
	if err != nil {
		panic(fmt.Sprintf("Failed to compile WASM module: %v", err))
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
	instance, err := r.InstantiateModule(ctx, compiled, config)
	if err != nil {
		panic(fmt.Sprintf("Failed to instantiate WASM module: %v", err))
	}
	defer instance.Close(ctx)

	// Run the "_start" function, which is the standard entry point for WASI modules
	_, err = instance.ExportedFunction("_start").Call(ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to run WASM module: %v", err))
	}

	// Read and print the output from the WASM module
	fmt.Printf("Output from WASM: %s", stdout.String())
}
