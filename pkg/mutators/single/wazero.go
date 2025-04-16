//go:build plugins
// +build plugins

package mutators

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/batmac/ccat/pkg/log"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"
	// "github.com/tetratelabs/wazero/api"
)

func init() {
	singleRegister("wasm", wasm, withDescription("a wasi (wasm) module to apply (path as first argument)"),
		withConfigBuilder(stdConfigString),
		withCategory("plugin"),
	)
}

func wasm(w io.WriteCloser, r io.ReadCloser, arg any) (int64, error) {
	ctx := context.Background()
	moduleFile := arg.(string)

	log.Debugln("Create a new WebAssembly runtime")
	runtime := wazero.NewRuntime(ctx)
	defer runtime.Close(ctx)

	log.Debugln("Instantiate WASI in the runtime")
	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

	log.Debugln("Load your WASM module")
	wasmBytes, err := os.ReadFile(moduleFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to read WASM module: %v", err))
	}

	log.Debugln("Define the configuration with the in-memory file system")
	config := wazero.NewModuleConfig().
		WithStdin(r).
		WithStdout(w).
		WithStderr(os.Stderr) // Set stderr to the standard error of the host process

	log.Debugln("Instantiate the WASM module with the configured input/output")
	if _, err := runtime.InstantiateWithConfig(ctx, wasmBytes, config); err != nil {
		var exitErr *sys.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() != 0 {
			log.Debugf("exit_code: %d\n", exitErr.ExitCode())
		} else {
			log.Fatal(err)
		}
	}
	w.Close()

	return 0, nil
}
