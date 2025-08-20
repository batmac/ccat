//go:build !fileonly
// +build !fileonly

package openers

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/batmac/ccat/pkg/globalctx"
	"github.com/batmac/ccat/pkg/log"
)

var (
	sseOpenerName        = "sse"
	sseOpenerDescription = "stream Server-Sent Events via sse://"
)

type sseOpener struct {
	name, description string
}

func init() {
	register(&sseOpener{
		name:        sseOpenerName,
		description: sseOpenerDescription,
	})
}

func (f sseOpener) Name() string {
	return f.name
}

func (f sseOpener) Description() string {
	return f.description
}

func (f sseOpener) Open(s string, _ bool) (io.ReadCloser, error) {
	// Convert sse:// to https://
	url := strings.Replace(s, "sse://", "https://", 1)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	// Set SSE headers
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("SSE request failed with status: %d", resp.StatusCode)
	}
	
	// Mark as streaming content to bypass highlighter
	globalctx.Set("hintSSEStream", true)
	
	// Return a custom ReadCloser that parses SSE events
	return &sseReadCloser{
		resp:    resp,
		scanner: bufio.NewScanner(resp.Body),
		buf:     &bytes.Buffer{},
	}, nil
}

func (f sseOpener) Evaluate(s string) float32 {
	if strings.HasPrefix(s, "sse://") {
		return 0.9
	}
	return 0
}

// sseReadCloser wraps the HTTP response and parses SSE events
type sseReadCloser struct {
	resp    *http.Response
	scanner *bufio.Scanner
	buf     *bytes.Buffer
	closed  bool
}

func (s *sseReadCloser) Read(p []byte) (int, error) {
	if s.closed {
		return 0, io.EOF
	}
	
	// If buffer has data, return it first
	if s.buf.Len() > 0 {
		return s.buf.Read(p)
	}
	
	// Parse next SSE event
	for s.scanner.Scan() {
		line := s.scanner.Text()
		
		// Skip empty lines
		if line == "" {
			continue
		}
		
		// Only output data lines
		if strings.HasPrefix(line, "data:") {
			eventData := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			s.buf.WriteString(eventData)
			s.buf.WriteString("\n")
			return s.buf.Read(p)
		}
		// Skip all other SSE fields (event:, id:, retry:, comments)
	}
	
	// Check for scanner errors
	if err := s.scanner.Err(); err != nil {
		log.Debugf("SSE scanner error: %v", err)
		return 0, err
	}
	
	// End of stream
	return 0, io.EOF
}

func (s *sseReadCloser) Close() error {
	if s.closed {
		return nil
	}
	s.closed = true
	return s.resp.Body.Close()
}