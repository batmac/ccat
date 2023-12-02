package miniclaude

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"braces.dev/errtrace"
	"github.com/batmac/ccat/pkg/log"
)

const (
	PromptHuman = "\n\nHuman:"
	PromptAI    = "\n\nAssistant:"

	// https://docs.anthropic.com/claude/reference/selecting-a-model
	ModelClaudeLatest        = "claude-2"           // latest model family, manually updated
	ModelClaude2             = "claude-2"           // latest major version
	ModelClaude21            = "claude-2.1"         // latest full version
	ModelClaudeInstantLatest = "claude-instant-1"   // latest instant model family, manually updated
	ModelClaudeInstant1      = "claude-instant-1"   // latest instant major version
	ModelClaudeInstant12     = "claude-instant-1.2" // latest instant full version
)

var (
	BaseURL                  = "https://api.anthropic.com"
	APIVersion               = "2023-06-01"
	DefaultMaxTokensToSample = 1200
	DefaultStopSequences     = []string{PromptHuman}
)

//nolint:tagliatelle
type Metadata struct {
	UserID string `json:"user_id"`
}

//nolint:tagliatelle
type SamplingParameters struct {
	Temperature       *float64          `json:"temperature,omitempty"`
	TopK              *int              `json:"top_k,omitempty"`
	TopP              *float64          `json:"top_p,omitempty"`
	Tags              map[string]string `json:"tags,omitempty"`
	Prompt            string            `json:"prompt"`
	Model             string            `json:"model"`
	StopSequences     []string          `json:"stop_sequences"`
	MaxTokensToSample int               `json:"max_tokens_to_sample"`
	Stream            bool              `json:"stream"`
	Metadata          Metadata          `json:"metadata,omitempty"`
}

//nolint:tagliatelle
type response struct {
	Completion string `json:"completion"`
	StopReason string `json:"stop_reason"`
	Model      string `json:"model"`
	Stop       string `json:"stop"`
	LogID      string `json:"log_id"`
	Exception  string `json:"exception"`
	Truncated  bool   `json:"truncated"`
}

func WrapPrompt(human, ai string) string {
	if len(ai) == 0 {
		return fmt.Sprintf("%s %s%s", PromptHuman, human, PromptAI)
	}
	return fmt.Sprintf("%s %s%s %s", PromptHuman, human, PromptAI, ai)
}

func NewSimpleSamplingParameters(prompt string, model string) *SamplingParameters {
	return &SamplingParameters{
		Prompt:            WrapPrompt(prompt, ""),
		MaxTokensToSample: DefaultMaxTokensToSample,
		StopSequences:     DefaultStopSequences,
		Model:             model,
		Stream:            true,
		Metadata:          Metadata{UserID: "ccat"},
	}
}

type Client struct {
	HTTPClient *http.Client
	C          chan string
	Endpoint   string
	APIKey     string
	APIVersion string
}

func New() *Client {
	return &Client{
		Endpoint:   BaseURL + "/v1/complete",
		APIKey:     os.Getenv("ANTHROPIC_API_KEY"),
		APIVersion: APIVersion,
		HTTPClient: &http.Client{},
		C:          make(chan string, 5), // buffer because we don't want to block the stream
	}
}

func (c *Client) Stream(sp *SamplingParameters) error {
	sp.Stream = true
	if sp.Prompt == "" {
		return errtrace.Wrap(fmt.Errorf("prompt is required"))
	}
	if sp.Model == "" {
		sp.Model = ModelClaudeLatest
	}

	data, err := json.Marshal(sp)
	if err != nil {
		return errtrace.Wrap(err)
	}
	req, err := http.NewRequest(http.MethodPost, c.Endpoint, bytes.NewReader(data))
	if err != nil {
		return errtrace.Wrap(err)
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("anthropic-version", c.APIVersion)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return errtrace.Wrap(err)
	}
	defer resp.Body.Close()
	defer close(c.C)

	if resp.StatusCode != http.StatusOK {
		log.Printf("http unexpected status code: %s", resp.Status)
		return errtrace.Wrap(fmt.Errorf("http unexpected status code: %d (%s)", resp.StatusCode, resp.Status))
	}

	scanner := bufio.NewScanner(resp.Body)
	for {
		scanner.Scan()
		eventName := strings.TrimPrefix(scanner.Text(), "event:")
		eventName = strings.TrimSpace(eventName)
		scanner.Scan()
		eventData := strings.TrimPrefix(scanner.Text(), "data:")
		scanner.Scan()
		if len(scanner.Text()) != 0 {
			return errtrace.Wrap(fmt.Errorf("unexpected line: %s", scanner.Text()))
		}
		if scanner.Err() != nil {
			return errtrace.Wrap(fmt.Errorf("error reading from scanner: %w", scanner.Err()))
		}
		log.Debugf("received event '%s': %s", eventName, eventData)
		if eventName == "error" {
			log.Printf("error: %s", eventData)
		}
		if eventName != "completion" && eventName != "ping" {
			return errtrace.Wrap(fmt.Errorf("unexpected event name '%s'", eventName))
		}

		var r response
		err = json.Unmarshal([]byte(eventData), &r)
		if err != nil {
			log.Debugf("error unmarshalling: %s\n%s\n", err, eventData)
			return errtrace.Wrap(err)
		}
		if r.StopReason != "" || r.Exception != "" {
			log.Debugf("stop reason: %s, stop: %s, exception: %s \n", r.StopReason, strconv.Quote(r.Stop), r.Exception)
			c.C <- "\n"
			break
		}

		if eventName == "completion" {
			c.C <- r.Completion
		}
	}
	return nil
}
