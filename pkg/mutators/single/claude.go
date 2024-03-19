package mutators

import (
	"io"
	"strconv"

	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/miniclaude"
	"github.com/batmac/ccat/pkg/secretprovider"
)

func init() {
	singleRegister("claude", claude,
		withDescription("ask Anthropic Claude, X:<unlimited> max replied tokens, optional second arg is the model, optional third arg is the preprompt (needs a valid key in $ANTHROPIC_API_KEY)"),
		withConfigBuilder(stdConfigStrings(0, 3)),
		withHintSlow(), // output asap (when no other mutator is used),
		withCategory("external APIs"),
	)
}

func claude(w io.WriteCloser, r io.ReadCloser, conf any) (int64, error) {
	args := conf.([]string)
	maxTokens := 1000
	var err error
	if len(args) > 0 && args[0] != "" {
		maxTokens, err = strconv.Atoi(args[0])
		if err != nil {
			log.Println("first arg: ", err)
		}
	}

	model := miniclaude.ModelClaude3Haiku
	if len(args) >= 2 && args[1] != "" {
		model = args[1]
	}

	prePrompt := ""
	if len(args) >= 3 && args[2] != "" {
		prePrompt = args[2] + ":\n"
	}

	log.Debugln("model: ", model)
	log.Debugln("maxTokens: ", maxTokens)
	log.Debugln("prePrompt (system): ", prePrompt)
	key, _ := secretprovider.GetSecret("anthropic", "ANTHROPIC_API_KEY")
	if key == "" {
		log.Fatal("ANTHROPIC_API_KEY environment variable is not set")
	}

	prompt, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}

	mr := &miniclaude.MessagesRequest{
		Model: model,
		Messages: []miniclaude.Message{
			{
				Role: miniclaude.RoleUser,
				Content: []miniclaude.ContentBlock{
					{
						Type: miniclaude.ContentTypeText,
						Text: string(prompt),
					},
				},
			},
		},
		MaxTokens: maxTokens,
	}

	if prePrompt != "" {
		mr.System = prePrompt
	}

	request := miniclaude.NewMessagesRequest()
	request.APIKey = key

	go func() {
		if key == "CI" {
			log.Println("ANTHROPIC_API_KEY is set to CI, using fake response")
			request.C <- "fake"
			request.C <- ""
			close(request.C)
			return
		}

		err := request.Stream(mr)
		if err != nil {
			log.Println("request.Stream: ", err)
		}
	}()

	var total int64
	for s := range request.C {
		if s == "" {
			log.Debugln("empty string")
		}
		n, _ := w.Write([]byte(s))
		total += int64(n)
	}
	return total, nil
}
