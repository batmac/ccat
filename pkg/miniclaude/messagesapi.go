package miniclaude

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/batmac/ccat/pkg/log"
)

const (
	ModelClaude3Opus   = "claude-3-opus-20240229"
	ModelClaude3Sonnet = "claude-3-sonnet-20240229"
	ModelClaude3Haiku  = "claude-3-haiku-20240307"

	// legacy models
	ModelClaude21        = "claude-2.1"
	ModelClaudeInstant12 = "claude-instant-1.2"
)

var (
	ModelAliases = map[string]string{
		"opus":   ModelClaude3Opus,
		"sonnet": ModelClaude3Sonnet,
		"haiku":  ModelClaude3Haiku,
	}

	BaseURL                  = "https://api.anthropic.com"
	APIVersion               = "2023-06-01"
	DefaultMaxTokensToSample = 1200
)

func NewMessagesRequest() *Request {
	return &Request{
		Endpoint:   BaseURL + "/v1/messages",
		APIKey:     os.Getenv("ANTHROPIC_API_KEY"),
		APIVersion: APIVersion,
		HTTPClient: &http.Client{},
		C:          make(chan string, 5), // buffer because we don't want to block the stream
	}
}

func (c *Request) Stream(mr *MessagesRequest) error {
	mr.Stream = true

	if m, ok := ModelAliases[strings.ToLower(mr.Model)]; ok {
		mr.Model = m
	}
	if mr.Model == "" {
		mr.Model = ModelClaude3Haiku
	}
	if mr.MaxTokens == 0 {
		mr.MaxTokens = DefaultMaxTokensToSample
	}

	data, err := json.Marshal(mr)
	if err != nil {
		return err
	}
	log.Debugf("sending request: %s", data)
	req, err := http.NewRequest(http.MethodPost, c.Endpoint, bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("anthropic-version", c.APIVersion)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	defer close(c.C)

	if resp.StatusCode != http.StatusOK {
		log.Printf("http unexpected status code: %s", resp.Status)
		// dump the body
		io.Copy(os.Stderr, resp.Body)
		return fmt.Errorf("http unexpected status code: %d (%s)", resp.StatusCode, resp.Status)
	}

	scanner := bufio.NewScanner(resp.Body)
stream:
	for {
		scanner.Scan()
		eventName := strings.TrimPrefix(scanner.Text(), "event:")
		eventName = strings.TrimSpace(eventName)
		scanner.Scan()
		eventData := strings.TrimPrefix(scanner.Text(), "data:")
		scanner.Scan() // empty line
		if len(scanner.Text()) != 0 {
			return fmt.Errorf("unexpected line: %s", scanner.Text())
		}
		if scanner.Err() != nil {
			return fmt.Errorf("error reading from scanner: %w", scanner.Err())
		}
		log.Debugf("received event '%s': %s", eventName, eventData)
		if eventName == "error" {
			log.Printf("error: %s", eventData)
		}
		var d map[string]interface{}
		err = json.Unmarshal([]byte(eventData), &d)
		if err != nil {
			log.Debugf("error unmarshalling: %s\n%s\n", err, eventData)
			return err
		}

		switch eventName {
		case "ping":
			// no-op
		case "message_start":
			// no-op
		case "content_block_start":
			c.C <- d["content_block"].(map[string]any)["text"].(string)
		case "content_block_delta":
			c.C <- d["delta"].(map[string]any)["text"].(string)
		case "content_block_stop":
			// no-op
		case "message_delta":
			// no-op
		case "message_stop":
			break stream
		case "error":
			break stream
		default:
			log.Println("unexpected event name: ", eventName, eventData)
		}
	}
	log.Debugln("stream ended")
	return nil
}
