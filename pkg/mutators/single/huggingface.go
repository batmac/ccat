package mutators

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/utils"
)

// https://huggingface.co/docs/api-inference/detailed_parameters

// Supported tasks:

// Fill Mask task
// Tries to fill in a hole with a missing word (token to be precise). That’s the base task for BERT models.
// Recommended model: bert-base-uncased (it’s a simple model, but fun to play with).

// Summarization task
// This task is well known to summarize longer text into shorter text. Be careful, some models have a maximum length of input.
// That means that the summary cannot handle full books for instance. Be careful when choosing your model.
// If you want to discuss your summarization needs, please get in touch with us: <api-enterprise@huggingface.co>
// Recommended model: facebook/bart-large-cnn.

// Text Classification task
// Usually used for sentiment-analysis this will output the likelihood of classes of an input.
// Recommended model: distilbert-base-uncased-finetuned-sst-2-english

// Text Generation task
// Use to continue text from a prompt. This is a very generic task.
// Recommended model: gpt2 (it’s a simple model, but fun to play with).

// Token Classification task
// Usually used for sentence parsing, either grammatical, or Named Entity Recognition (NER) to understand keywords contained within text.
// Recommended model: dbmdz/bert-large-cased-finetuned-conll03-english

// Translation task
// This task is well known to translate text from one language to another
// Recommended model: Helsinki-NLP/opus-mt-ru-en. Helsinki-NLP uploaded many models with many language pairs.
// Recommended model: t5-base.

var HuggingFaceCommonTasks = map[string]string{
	"fillmask":        "bert-base-uncased",
	"summarization":   "facebook/bart-large-cnn",
	"classification":  "distilbert-base-uncased-finetuned-sst-2-english",
	"text-generation": "gpt2",
	"ner":             "dbmdz/bert-large-cased-finetuned-conll03-english",
	"translation":     "t5-base",
	"bloom":           "bigscience/bloom",
	"bloomz":          "bigscience/bloomz",
}

type HuggingFaceRequest struct {
	Inputs  string         `json:"inputs"`
	Options map[string]any `json:"options"`
}

func init() {
	singleRegister("huggingface", huggingface,
		withDescription("ask HuggingFace for simple tasks, optional arg is the model (needs a valid key in $HUGGING_FACE_HUB_TOKEN)"),
		withConfigBuilder(stdConfigStrings(0, 1)),
		withAliases("hf"),
		withCategory("external APIs"),
	)
}

func huggingface(w io.WriteCloser, r io.ReadCloser, conf any) (int64, error) {
	arg := conf.([]string)
	baseURL := "https://api-inference.huggingface.co/models/"

	token, source, err := getHuggingFaceToken()
	if err != nil && os.Getenv("CI") != "CI" {
		return 0, err
	}
	model := "bigscience/bloom"
	if len(arg) >= 1 && arg[0] != "" {
		model = arg[0]
	}

	log.Debugf("task aliases: %v\n", HuggingFaceCommonTasks)
	if m, ok := HuggingFaceCommonTasks[model]; ok {
		model = m
	}

	log.Debugln("token: from ", source)
	log.Debugln("model: ", model)

	input, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}
	request, err := json.Marshal(HuggingFaceRequest{Inputs: string(input), Options: map[string]any{"wait_for_model": true}})
	if err != nil {
		return 0, err
	}

	log.Debugf("request: %s\n", request)

	req, err := http.NewRequest(http.MethodPost, baseURL+model, bytes.NewReader(request))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "ccat")

	if os.Getenv("CI") == "CI" {
		_, _ = w.Write([]byte("fake"))
		return 0, nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("error: %s", resp.Status)
	}

	got, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	n, err := w.Write(got)
	return int64(n), err
}

//nolint:gosec
func getHuggingFaceToken() (string, string, error) {
	// HUGGING_FACE_HUB_TOKEN,
	// then HF_API_KEY,
	// then the content of the file $HF_HOME/token
	// then the content of the file ~/.huggingface/token
	// and finally the content of the file ~/.cache/huggingface/token
	token, source := os.Getenv("HUGGING_FACE_HUB_TOKEN"), "HUGGING_FACE_HUB_TOKEN"
	if token == "" {
		token, source = os.Getenv("HF_API_KEY"), "HF_API_KEY"
	}
	for _, path := range []string{"$HF_HOME/token", "~/.huggingface/token", "~/.cache/huggingface/token"} {
		if token != "" {
			break
		}
		content, _ := os.ReadFile(utils.ExpandPath(path))
		token, source = string(content), path
	}

	if token == "" || os.Getenv("CI") == "CI" {
		return "", "", fmt.Errorf("no HuggingFace token found")
	}
	return token, source, nil
}
