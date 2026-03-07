//go:build !fileonly
// +build !fileonly

package openers

import (
	"fmt"
	"io"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/term"
	"github.com/batmac/ccat/pkg/utils"
)

var (
	clipboardOpenerName        = "clipboard"
	clipboardOpenerDescription = "get content from the system clipboard via cb://"
)

type clipboardOpener struct {
	name, description string
}

func init() {
	register(&clipboardOpener{
		name:        clipboardOpenerName,
		description: clipboardOpenerDescription,
	})
}

func (c clipboardOpener) Name() string {
	return c.name
}

func (c clipboardOpener) Description() string {
	return c.description
}

func (c clipboardOpener) Open(s string, _ bool) (io.ReadCloser, error) {
	// Verify protocol
	protocol, _, found := strings.Cut(s, "://")
	if !found || protocol != "cb" {
		return nil, fmt.Errorf("clipboard opener requires cb:// protocol")
	}

	// Warn if in SSH or container - reading from remote clipboard
	if term.IsSSH() {
		log.Println("WARNING: cb:// reads from remote server clipboard in SSH session")
		log.Println("         Consider using stdin or clipboard forwarding for local clipboard access")
	} else if utils.IsRunningInContainer() {
		log.Println("WARNING: cb:// reads from container clipboard, not host clipboard")
	}

	log.Debugln("Reading from clipboard...")

	// Read from clipboard
	content, err := clipboard.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read from clipboard: %w", err)
	}

	log.Debugf("Read %d bytes from clipboard\n", len(content))

	// Return content as a ReadCloser
	return io.NopCloser(strings.NewReader(content)), nil
}

func (c clipboardOpener) Evaluate(s string) float32 {
	if strings.HasPrefix(s, "cb://") {
		return 1.0 // High priority for cb:// URLs
	}
	return 0
}
