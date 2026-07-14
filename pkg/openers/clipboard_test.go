//go:build !fileonly
// +build !fileonly

package openers

import (
	"io"
	"strings"
	"testing"

	"github.com/atotto/clipboard"
)

func TestClipboardOpener_Evaluate(t *testing.T) {
	opener := &clipboardOpener{
		name:        clipboardOpenerName,
		description: clipboardOpenerDescription,
	}

	tests := []struct {
		name     string
		input    string
		expected float32
	}{
		{"valid cb URL", "cb://", 1.0},
		{"valid cb URL with content", "cb://anything", 1.0},
		{"http URL", "http://example.com", 0},
		{"file path", "/path/to/file", 0},
		{"echo URL", "echo://test", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := opener.Evaluate(tt.input)
			if score != tt.expected {
				t.Errorf("Evaluate(%q) = %v, want %v", tt.input, score, tt.expected)
			}
		})
	}
}

func TestClipboardOpener_Open(t *testing.T) {
	opener := &clipboardOpener{
		name:        clipboardOpenerName,
		description: clipboardOpenerDescription,
	}

	// Test 1: Write to clipboard and read it back
	testContent := "test clipboard content"
	err := clipboard.WriteAll(testContent)
	if err != nil {
		t.Skipf("Clipboard not available (might be in CI): %v", err)
	}

	rc, err := opener.Open("cb://", false)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	if string(content) != testContent {
		t.Errorf("Open() content = %q, want %q", string(content), testContent)
	}

	// Test 2: Invalid protocol
	_, err = opener.Open("http://example.com", false)
	if err == nil {
		t.Error("Open() with invalid protocol should return error")
	}
	if !strings.Contains(err.Error(), "cb://") {
		t.Errorf("Open() error should mention cb:// protocol, got: %v", err)
	}
}

func TestClipboardOpener_Name(t *testing.T) {
	opener := &clipboardOpener{
		name:        clipboardOpenerName,
		description: clipboardOpenerDescription,
	}

	if name := opener.Name(); name != clipboardOpenerName {
		t.Errorf("Name() = %q, want %q", name, clipboardOpenerName)
	}
}

func TestClipboardOpener_Description(t *testing.T) {
	opener := &clipboardOpener{
		name:        clipboardOpenerName,
		description: clipboardOpenerDescription,
	}

	if desc := opener.Description(); desc != clipboardOpenerDescription {
		t.Errorf("Description() = %q, want %q", desc, clipboardOpenerDescription)
	}
}
