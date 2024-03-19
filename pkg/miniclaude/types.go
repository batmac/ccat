package miniclaude

import (
	"net/http"
)

type Request struct {
	HTTPClient *http.Client
	C          chan string
	Endpoint   string
	APIKey     string
	APIVersion string
}

type MessagesResponse struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"` // always "message"
	Error        ErrorBlock     `json:"error,omitempty"`
	Role         string         `json:"role"` // always "assistant"
	Content      []ContentBlock `json:"content"`
	Model        string         `json:"model"`
	StopReason   string         `json:"stop_reason"`
	StopSequence string         `json:"stop_sequence,omitempty"`
	Usage        Usage          `json:"usage"`
}

const (
	StopReasonMaxTokens    string = "max_tokens"
	StopReasonStopSequence string = "stop_sequence"
	StopReasonEndTurn      string = "end_turn"
)

type ErrorBlock struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
	// Source is only present if the type is "image"
	Source *ContentBlockSource `json:"source,omitempty"`
}

const (
	ContentTypeText  string = "text"
	ContentTypeImage string = "image"
)

type ContentBlockSource struct {
	Type      string `json:"type"` // always "base64"
	Data      string `json:"data"`
	MediaType string `json:"media_type"`
}

const (
	ContentBlockMediaTypeImagePng  string = "image/png"
	ContentBlockMediaTypeImageJpeg string = "image/jpeg"
	ContentBlockMediaTypeImageGif  string = "image/gif"
	ContentBlockMediaTypeImageWebp string = "image/webp"
)

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type MessagesRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	System    string    `json:"system,omitempty"`
	MaxTokens int       `json:"max_tokens"`
	Stream    bool      `json:"stream,omitempty"`
	// Metadata  MetadataObject `json:"metadata,omitempty"`
}

type Message struct {
	Role    string         `json:"role"`
	Content []ContentBlock `json:"content"`
}

const (
	RoleUser      string = "user"
	RoleAssistant string = "assistant"
)

// type MetadataObject struct {
// StopSequences []string `json:"stop_sequences,omitempty"`
// Stream        bool     `json:"stream,omitempty"`
// Temperature   *float64 `json:"temperature,omitempty"`
// TopP          float64  `json:"top_p,omitempty"`
// TopK          int      `json:"top_k,omitempty"`
// }
