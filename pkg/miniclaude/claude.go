package miniclaude

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/batmac/ccat/pkg/log"
)

const (
	PromptHuman = "\n\nHuman:"
	PromptAI    = "\n\nAssistant:"
	MessageDone = "[DONE]"

	ModelClaudeV1         string = "claude-v1"
	ModelClaudeV10        string = "claude-v1.0"
	ModelClaudeV12        string = "claude-v1.2"
	ModelClaudeInstantV1  string = "claude-instant-v1"
	ModelClaudeInstantV10 string = "claude-instant-v1.0"
)

var (
	BaseURL                  = "https://api.anthropic.com"
	DefaultMaxTokensToSample = 1200
	DefaultStopSequences     = []string{PromptHuman}
)

//nolint:tagliatelle
type SamplingParameters struct {
	Prompt            string            `json:"prompt"`
	MaxTokensToSample int               `json:"max_tokens_to_sample"`
	StopSequences     []string          `json:"stop_sequences"`
	Model             string            `json:"model"`
	Stream            bool              `json:"stream"`
	Temperature       *float64          `json:"temperature,omitempty"`
	TopK              *int              `json:"top_k,omitempty"`
	TopP              *float64          `json:"top_p,omitempty"`
	Tags              map[string]string `json:"tags,omitempty"`
}

//nolint:tagliatelle
type response struct {
	Completion string `json:"completion"`
	Stop       string `json:"stop"`
	StopReason string `json:"stop_reason"`
	Truncated  bool   `json:"truncated"`
	LogID      string `json:"log_id"`
	Model      string `json:"model"`
	Exception  string `json:"exception"`
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
	}
}

type Client struct {
	Endpoint   string
	APIKey     string
	HTTPClient *http.Client
	C          chan string
}

func New() *Client {
	return &Client{
		Endpoint:   BaseURL + "/v1/complete",
		APIKey:     os.Getenv("ANTHROPIC_API_KEY"),
		HTTPClient: &http.Client{},
		C:          make(chan string, 5), // buffer because we don't want to block the stream
	}
}

func (c *Client) Stream(sp *SamplingParameters) error {
	sp.Stream = true
	if sp.Prompt == "" {
		return fmt.Errorf("prompt is required")
	}
	if sp.Model == "" {
		sp.Model = ModelClaudeV12
	}

	data, err := json.Marshal(sp)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.Endpoint, bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	defer close(c.C)

	var previousResponse response
	scanner := bufio.NewScanner(resp.Body)
	event := bytes.NewBuffer(nil)
	for scanner.Scan() {
		line := scanner.Text()
		// if empty, we have an event in "event"
		if len(line) == 0 {
			log.Debugf("event: %s\n", event.String())
			var r response
			err = json.Unmarshal(event.Bytes(), &r)
			event.Reset()
			if err != nil {
				log.Debugf("error unmarshalling: %s\n%s\n", err, event.String())
				continue
			}
			c.C <- strings.TrimPrefix(r.Completion, previousResponse.Completion)
			previousResponse = r
			continue
		}
		line = strings.TrimPrefix(line, "data: ")
		if line == MessageDone {
			log.Debugf("done: %s\n", event.String()+line)
			if previousResponse.StopReason != "stop_sequence" && previousResponse.StopReason != "max_tokens" {
				log.Printf("unexpected stop reason: %s", previousResponse.StopReason)
			}
			break
		}
		event.WriteString(line)
	}
	if err := scanner.Err(); err != nil {
		log.Println("reading standard input:", err)
	}
	return nil
}
