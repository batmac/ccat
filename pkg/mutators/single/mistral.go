package mutators

import (
	"io"
	"strconv"
	"strings"

	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/secretprovider"

	"github.com/gage-technologies/mistral-go"
)

// https://platform.openai.com/docs/guides/chat

func init() {
	singleRegister("mistralai", mistralai,
		withDescription("ask MistralAI, X:<unlimited> max replied tokens, the optional second arg is the model (Requires a valid key in $MISTRAL_API_KEY)"),
		withConfigBuilder(stdConfigStrings(0, 2)),
		withAliases("mistral"),
		withHintSlow(), // output asap (when no other mutator is used)
		withCategory("external APIs"),
	)
}

func mistralai(w io.WriteCloser, r io.ReadCloser, conf any) (int64, error) {
	args := conf.([]string)
	model := "mistral-tiny"
	maxTokens := 4000
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

	key, _ := secretprovider.GetSecret("mistralai", "MISTRAL_API_KEY")
	if key == "" {
		log.Fatal("MISTRAL_API_KEY environment variable is not set")
	}

	log.Debugln("model: ", model)
	log.Debugln("maxTokens: ", maxTokens)

	client := mistral.NewMistralClientDefault(key)
	// log.Debugf("models: %+v", listModels(client))

	prompt, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}

	req := []mistral.ChatMessage{{Content: string(prompt), Role: mistral.RoleUser}}
	log.Debugf("request: %#v", req)
	if key == "CI" {
		log.Println("MISTRAL_API_KEY is set to CI, using fake response")
		return io.Copy(w, strings.NewReader("CI"))
	}
	params := mistral.DefaultChatRequestParams
	params.MaxTokens = maxTokens
	stream, err := client.ChatStream(model, req, &params)
	if err != nil {
		return 0, err
	}

	defer func() {
		if _, err = w.Write([]byte("\n")); err != nil {
			log.Println(err)
		}
	}()

	var totalWritten int64
	var steps int
	for chunk := range stream {
		if chunk.Error != nil {
			return 0, chunk.Error
		}
		log.Debugf("chunk: %#v", chunk)
		n, err := w.Write([]byte(chunk.Choices[0].Delta.Content))
		if err != nil {
			return 0, err
		}
		totalWritten += int64(n)
		steps++
	}
	log.Debugf("finished after %d steps.", steps)
	return totalWritten, nil
}
