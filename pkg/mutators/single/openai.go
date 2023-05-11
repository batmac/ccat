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

func init() {
	singleRegister("chatgpt", chatgpt,
		withDescription("ask OpenAI ChatGPT, X:<unlimited> max replied tokens, the optional second arg is the model (Requires a valid key in $OPENAI_API_KEY, optional custom endpoint in $OPENAI_BASE_URL.)"),
		withConfigBuilder(stdConfigStrings(0, 3)),
		withAliases("cgpt"),
		withHintSlow(), // output asap (when no other mutator is used)
		withCategory("external APIs"),
	)
}

func chatgpt(w io.WriteCloser, r io.ReadCloser, conf any) (int64, error) {
	args := conf.([]string)
	model := gpt.GPT3Dot5Turbo
	maxTokens := 0 // unlimited
	var err error
	if len(args) > 0 && args[0] != "" {
		maxTokens, err = strconv.Atoi(args[0])
		if err != nil {
			log.Println("first arg: ", err)
		}
	}
	if len(args) >= 2 && args[1] != "" {
		model = args[1]
	}

	prePrompt := ""
	if len(args) >= 3 && args[2] != "" {
		prePrompt = args[2] + ":\n"
	}

	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		log.Fatal("OPENAI_API_KEY environment variable is not set")
	}

	config := gpt.DefaultConfig(key)
	customBaseURL := os.Getenv("OPENAI_BASE_URL")
	if customBaseURL != "" {
		config.BaseURL = customBaseURL
	}

	log.Debugln("baseURL: ", config.BaseURL)

	log.Debugln("model: ", model)
	log.Debugln("maxTokens: ", maxTokens)
	log.Debugln("prePrompt: ", prePrompt)

	client := gpt.NewClientWithConfig(config)
	ctx := context.Background()
	// log.Debugf("models: %+v", listModels(client))

	prompt, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}

	req := gpt.ChatCompletionRequest{
		Model: model,
		Messages: []gpt.ChatCompletionMessage{
			{Role: "user", Content: prePrompt + string(prompt)},
		},
		MaxTokens:        maxTokens,
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
	if key == "CI" {
		log.Println("OPENAI_API_KEY is set to CI, using fake response")
		return io.Copy(w, strings.NewReader("CI"))
	}
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
