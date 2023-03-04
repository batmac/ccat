package mutators

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/batmac/ccat/pkg/log"
	gpt "github.com/sashabaranov/go-gpt3"
)

// https://platform.openai.com/docs/guides/chat

var defaultModel = gpt.GPT3Dot5Turbo

func init() {
	singleRegister("chatgpt", chatgpt,
		withDescription("ask OpenAI ChatGPT, X:4000 max replied tokens (needs a valid key in $OPENAI_API_KEY)"),
		withConfigBuilder(stdConfigUint64WithDefault(4000)),
	)
}

func chatgpt(w io.WriteCloser, r io.ReadCloser, conf any) (int64, error) {
	model := defaultModel
	maxTokens := conf.(uint64)
	log.Debugln("model: ", model)
	log.Debugln("maxTokens: ", maxTokens)
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		log.Fatal("OPENAI_API_KEY environment variable is not set")
	}

	client := gpt.NewClient(key)
	ctx := context.Background()

	prompt, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}

	req := gpt.ChatCompletionRequest{
		Model: model,
		Messages: []gpt.ChatCompletionMessage{
			{Role: "user", Content: string(prompt)},
		},
		MaxTokens:        int(maxTokens),
		Temperature:      0,
		TopP:             0,
		N:                0,
		Stream:           true,
		Stop:             []string{},
		PresencePenalty:  0,
		FrequencyPenalty: 0,
		LogitBias:        map[string]int{},
		User:             "",
	}
	log.Debugf("request: %#v", req)
	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return 0, err
	}
	defer stream.Close()
	//nolint:bodyclose // body is closed in stream.Close()
	if stream.GetResponse().StatusCode != http.StatusOK && key != "CI" {
		return 0, errors.New(stream.GetResponse().Status)
	}

	defer func() {
		if _, err = w.Write([]byte("\n")); err != nil {
			log.Println(err)
		}
	}()

	var totalWritten int64
	var steps int
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			log.Debugf("Stream finished after %d steps, response=%#v", steps, response)
			return totalWritten, nil
		}

		if err != nil {
			log.Printf("Stream error after %d steps: %v\n", steps, err)
			return totalWritten, err
		}

		log.Debugf("%#v\n", response)
		output := response.Choices[0].Delta.Content
		n, err := w.Write([]byte(output))
		if err != nil {
			return totalWritten, err
		}
		totalWritten += int64(n)
		steps++
	}
}
