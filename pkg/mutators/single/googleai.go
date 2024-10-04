package mutators

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/batmac/ccat/pkg/log"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func init() {
	singleRegister("googleai", googleai,
		withDescription("googleai, X:gemini-1.5-flash is the model (Requires a valid key in $GOOGLE_API_KEY)"),
		withAliases("gai"),
		withHintSlow(), // output asap (when no other mutator is used)
		withCategory("external APIs"),
		withConfigBuilder(stdConfigStringWithDefault("gemini-1.5-flash")),
	)
}

func googleai(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	modelname := config.(string)

	ctx := context.Background()

	apikey := os.Getenv("GOOGLE_API_KEY")
	log.Debugln("masked apikey: ", mask(apikey))
	client, err := genai.NewClient(ctx, option.WithAPIKey(apikey))
	if err != nil {
		return 0, err
	}
	defer client.Close()

	model := client.GenerativeModel(modelname)

	prompt, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}

	iter := model.GenerateContentStream(ctx, genai.Text(prompt))

	totalWritten := int64(0)
	for {
		resp, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			log.Debugln("error: ", err)
			return 0, err
		}

		p := ""
		for i, c := range resp.Candidates[0].Content.Parts {
			log.Debugln("part ", i, ": ", c)
			if t, ok := c.(genai.Text); ok {
				p += string(t)
			}
		}

		n, err := io.WriteString(w, p)
		if err != nil {
			return 0, err
		}
		totalWritten += int64(n)
	}
	return totalWritten, nil
}

func mask(s string) string {
	if len(s) <= 8 {
		return "********"
	}
	return "****..." + s[len(s)-2:]
}
