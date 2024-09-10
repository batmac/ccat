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
		if err != nil {
			break
		}
		writer.WriteString(input)
		writer.Flush()
	}
}
