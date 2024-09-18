//go:build ignore

// compile with:
// GOOS=wasip1 GOARCH=wasm go build -o go.wasm go.go
// or
// tinygo build -o tgo.wasm -target=wasi -scheduler=none -no-debug main.go

package main

import (
	"bufio"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	for {
		input, err := reader.ReadString('\n')
		writer.WriteString(input)
		if err != nil {
			break
		}
	}
	writer.Flush()
}
