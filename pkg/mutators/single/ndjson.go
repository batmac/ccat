package mutators

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func init() {
	singleRegister("ndjsonindent", NDJSONIndent, withDescription("pretty-print NDJSON/JSON Lines"),
		withHintLexer("JSON"),
		withCategory("convert"))
}

func NDJSONIndent(w io.WriteCloser, r io.ReadCloser, arg any) (int64, error) {
	// Parse indent level from argument, default to 2
	indent := 2
	if arg != nil {
		if s, ok := arg.(string); ok {
			if i, err := strconv.Atoi(s); err == nil && i >= 0 {
				indent = i
			}
		}
	}

	scanner := bufio.NewScanner(r)
	var totalBytes int64
	var indentStr string
	if indent > 0 {
		indentStr = strings.Repeat(" ", indent)
	}

	for scanner.Scan() {
		line := scanner.Text()
		totalBytes += int64(len(line)) + 1 // +1 for newline

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			if _, err := fmt.Fprintln(w); err != nil {
				return totalBytes, err
			}
			continue
		}

		// Parse and pretty-print the JSON line
		var jsonObj interface{}
		if err := json.Unmarshal([]byte(line), &jsonObj); err != nil {
			// If it's not valid JSON, write the line as-is
			if _, err := fmt.Fprintln(w, line); err != nil {
				return totalBytes, err
			}
			continue
		}

		// Pretty-print the JSON
		var buf bytes.Buffer
		if err := json.Indent(&buf, []byte(line), "", indentStr); err != nil {
			// If indenting fails, write the line as-is
			if _, err := fmt.Fprintln(w, line); err != nil {
				return totalBytes, err
			}
			continue
		}

		if _, err := fmt.Fprintln(w, buf.String()); err != nil {
			return totalBytes, err
		}
	}

	if err := scanner.Err(); err != nil {
		return totalBytes, err
	}

	return totalBytes, nil
}