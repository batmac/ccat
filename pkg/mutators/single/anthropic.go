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
	maxTokens := 0
	var err error
	if len(args) > 0 && args[0] != "" {
		maxTokens, err = strconv.Atoi(args[0])
		if err != nil {
			log.Println("first arg: ", err)
		}
	}

	model := miniclaude.ModelClaudeLatest // latest model
	if len(args) >= 2 && args[1] != "" {
		model = args[1]
	}

	prePrompt := ""
	if len(args) >= 3 && args[2] != "" {
		prePrompt = args[2] + ":\n"
	}

	log.Debugln("model: ", model)
	log.Debugln("maxTokens: ", maxTokens)
	log.Debugln("prePrompt: ", prePrompt)
	key, _ := secretprovider.GetSecret("anthropic", "ANTHROPIC_API_KEY")
	if key == "" {
		log.Fatal("ANTHROPIC_API_KEY environment variable is not set")
	}

	prompt, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}

	sp := miniclaude.NewSimpleSamplingParameters(prePrompt+string(prompt), model)
	if maxTokens > 0 {
		sp.MaxTokensToSample = maxTokens
	}
	client := miniclaude.New()
	client.APIKey = key

	go func() {
		if key == "CI" {
			log.Println("ANTHROPIC_API_KEY is set to CI, using fake response")
			client.C <- "fake"
			client.C <- ""
			close(client.C)
			return
		}

		err := client.Stream(sp)
		if err != nil {
			log.Println("client.Stream: ", err)
		}
	}()

	var total int64
	for s := range client.C {
		if s == "" {
			log.Debugln("empty string")
		}
		n, _ := w.Write([]byte(s))
		total += int64(n)
	}
	return total, nil
}
