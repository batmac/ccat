package mutators

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

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
	// keep the keys in lowercase
	"fillmask":        "bert-base-uncased",
	"summarization":   "facebook/bart-large-cnn",
	"classification":  "distilbert-base-uncased-finetuned-sst-2-english",
	"text-generation": "gpt2",
	"ner":             "dbmdz/bert-large-cased-finetuned-conll03-english",
	"translation":     "t5-base",
	"bloom":           "bigscience/bloom",
	"bloomz":          "bigscience/bloomz",
	"chat":            "OpenAssistant/oasst-sft-4-pythia-12b-epoch-3.5",
	"starcoder":       "bigcode/starcoder",
}

type HuggingFaceRequest struct {
	Options map[string]any `json:"options"`
	Inputs  string         `json:"inputs"`
}

func init() {
	singleRegister("huggingface", huggingface,
		withDescription("ask HuggingFace for simple tasks, optional arg is the model (needs a valid key in $HUGGING_FACE_HUB_TOKEN, set HUGGING_FACE_ENDPOINT to use an Inference API endpoint)"),
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

	if m, ok := HuggingFaceCommonTasks[strings.ToLower(model)]; ok {
		model = m
	}

	url := baseURL + model
	if os.Getenv("HUGGING_FACE_ENDPOINT") != "" {
		url = os.Getenv("HUGGING_FACE_ENDPOINT")
		if len(arg) >= 1 && arg[0] != "" {
			log.Println("warning: HUGGING_FACE_ENDPOINT is set, ignoring model argument")
		}
	}

	log.Debugln("token: from ", source)
	log.Debugln("model: ", model)
	log.Debugln("url: ", url)

	input, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}
	request, err := json.Marshal(HuggingFaceRequest{Inputs: string(input), Options: map[string]any{"wait_for_model": true}})
	if err != nil {
		return 0, err
	}

	log.Debugf("request: %s\n", request)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(request))
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

	out, err := getGeneratedTextFromJSON(got)
	if err == nil {
		got = out
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

func getGeneratedTextFromJSON(jsonBytes []byte) ([]byte, error) {
	type Response struct {
		GeneratedText string `json:"generated_text"` //nolint:tagliatelle
	}
	var data []Response

	log.Debugf("json: %s\n", jsonBytes)

	err := json.Unmarshal(jsonBytes, &data)
	if err != nil {
		log.Printf("error: %s\n", err)
		return nil, err
	}
	log.Debugf("data: %v\n", data)
	if len(data[0].GeneratedText) > 0 {
		return []byte(data[0].GeneratedText), nil
	}
	return nil, fmt.Errorf("no generated_text found in json response")
}
