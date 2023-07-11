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

	// https://docs.anthropic.com/claude/reference/selecting-a-model
	ModelClaudeLatest        = "claude-2"           // latest model family, manually updated
	ModelClaude2             = "claude-2"           // latest major version
	ModelClaude20            = "claude-2.0"         // latest full version
	ModelClaudeInstantLatest = "claude-instant-1"   // latest instant model family, manually updated
	ModelClaudeInstant1      = "claude-instant-1"   // latest instant major version
	ModelClaudeInstant11     = "claude-instant-1.1" // latest instant full version

	// old deprecated models, keeped for compatibility
	ModelClaudeV1_100K         = "claude-v1-100k"
	ModelClaudeV10             = "claude-v1.0"
	ModelClaudeV12             = "claude-v1.2"
	ModelClaudeV13             = "claude-v1.3"
	ModelClaudeV13_100K        = "claude-v1.3-100k"
	ModelClaudeInstantV1       = "claude-instant-v1"
	ModelClaudeInstantV1_100K  = "claude-instant-v1-100k"
	ModelClaudeInstantV10      = "claude-instant-v1.0"
	ModelClaudeInstantV11      = "claude-instant-v1.1"
	ModelClaudeInstantV11_100K = "claude-instant-v1.1-100k"
)

var (
	BaseURL                  = "https://api.anthropic.com"
	DefaultMaxTokensToSample = 1200
	DefaultStopSequences     = []string{PromptHuman}
)

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
}

//nolint:tagliatelle
type response struct {
	Completion string `json:"completion"`
	Stop       string `json:"stop"`
	StopReason string `json:"stop_reason"`
	LogID      string `json:"log_id"`
	Model      string `json:"model"`
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
	}
}

type Client struct {
	HTTPClient *http.Client
	C          chan string
	Endpoint   string
	APIKey     string
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

	if resp.StatusCode != http.StatusOK {
		log.Printf("http unexpected status code: %s", resp.Status)
		return fmt.Errorf("http unexpected status code: %d (%s)", resp.StatusCode, resp.Status)
	}

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
