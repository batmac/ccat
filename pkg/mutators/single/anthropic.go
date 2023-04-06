package mutators

import (
	"io"
	"os"
	"strconv"

	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/miniclaude"
)

func init() {
	singleRegister("claude", claude,
		withDescription("ask Anthropic Claude, X:<unlimited> max replied tokens, the optional second arg is the model (needs a valid key in $ANTHROPIC_API_KEY)"),
		withConfigBuilder(stdConfigStrings(0, 2)),
		withHintSlow(), // output asap (when no other mutator is used)
	)
}

func claude(w io.WriteCloser, r io.ReadCloser, conf any) (int64, error) {
	args := conf.([]string)
	maxTokens := uint64(defaultMaxTokens)
	var err error
	if len(args) > 0 && args[0] != "" {
		maxTokens, err = strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			log.Println("first arg: ", err)
		}
	}

	model := miniclaude.ModelClaudeV12
	if len(args) >= 2 {
		model = args[1]
	}

	log.Debugln("model: ", model)
	log.Debugln("maxTokens: ", maxTokens)
	key := os.Getenv("ANTHROPIC_API_KEY")
	if key == "" {
		log.Fatal("ANTHROPIC_API_KEY environment variable is not set")
	}
	/* if key == "CI" {
		log.Println("ANTHROPIC_API_KEY is set to CI, using fake response")
		return io.Copy(w, strings.NewReader("CI"))
	} */

	prompt, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}

	sp := miniclaude.NewSimpleSamplingParameters(string(prompt), model)
	if maxTokens > 0 {
		sp.MaxTokensToSample = int(maxTokens)
	}
	client := miniclaude.New()
	client.APIKey = key
	go func() {
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
