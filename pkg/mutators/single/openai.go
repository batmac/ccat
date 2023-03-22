package mutators

import (
	"context"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/batmac/ccat/pkg/log"
	gpt "github.com/sashabaranov/go-openai"
)

// https://platform.openai.com/docs/guides/chat

const (
	defaultChatModel = gpt.GPT3Dot5Turbo
	defaultMaxTokens = 0 // unlimited
)

func init() {
	singleRegister("chatgpt", chatgpt,
		withDescription("ask OpenAI ChatGPT, X:<unlimited> max replied tokens, the optional second arg is the model (needs a valid key in $OPENAI_API_KEY)"),
		withConfigBuilder(stdConfigStrings(0, 2)),
		withAliases("cgpt"),
	)
}

func chatgpt(w io.WriteCloser, r io.ReadCloser, conf any) (int64, error) {
	args := conf.([]string)
	model := defaultChatModel
	maxTokens := uint64(defaultMaxTokens)
	var err error
	if len(args) > 0 {
		maxTokens, err = strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			log.Println("first arg: ", err)
		}
	}
	if len(args) >= 2 {
		model = args[1]
	}

	log.Debugln("model: ", model)
	log.Debugln("maxTokens: ", maxTokens)
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		log.Fatal("OPENAI_API_KEY environment variable is not set")
	}
	if key == "CI" {
		log.Println("OPENAI_API_KEY is set to CI, using fake response")
		return io.Copy(w, strings.NewReader("CI"))
	}

	client := gpt.NewClient(key)
	ctx := context.Background()
	// log.Debugf("models: %+v", listModels(client))

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

/* func listModels(c *gpt.Client) string {
	models, err := c.ListModels(context.Background())
	if err != nil {
		log.Debugln("listModels(): ", err)
		return ""
	}
	// convert models to json string
	modelsJSON, err := json.Marshal(models)
	if err != nil {
		log.Debugln("listModels(): ", err)
		return ""
	}
	return string(modelsJSON)
}
*/
