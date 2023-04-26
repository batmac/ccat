package mutators

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/batmac/ccat/pkg/log"
)

type deeplTranslation struct {
	DetectedSourceLanguage string `json:"detected_source_language"` //nolint:tagliatelle
	Text                   string `json:"text"`
}

type deeplTranslateResponse struct {
	Translations []deeplTranslation `json:"translations"`
}

func init() {
	defaultLanguage := "en"
	if t := os.Getenv("TARGET_LANGUAGE"); t != "" {
		defaultLanguage = t
	}
	singleRegister("deepl", deepl,
		withDescription("translate to X:en or $TARGET_LANGUAGE with deepl (needs a valid key in $DEEPL_API_KEY)"),
		withConfigBuilder(stdConfigStringWithDefault(defaultLanguage)),
		withCategory("external APIs"),
	)
}

func deepl(w io.WriteCloser, r io.ReadCloser, conf any) (int64, error) {
	targetLanguage := conf.(string)
	text, err := io.ReadAll(r) // NOT streamable
	if err != nil {
		return 0, err
	}
	if len(text) == 0 {
		return 0, nil
	}

	key := os.Getenv("DEEPL_API_KEY")
	if key == "" {
		log.Fatal("no key found in $DEEPL_API_KEY")
	}

	baseURL := "https://api-free.deepl.com/v2/translate"
	// if apikey doesnt end with ":fx", use paid API
	if key[len(key)-3:] != ":fx" {
		baseURL = "https://api.deepl.com/v2/translate"
	}

	data := url.Values{}
	data.Set("text", string(text))
	data.Set("target_lang", targetLanguage)

	// Create request
	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "DeepL-Auth-Key "+key)

	// Send request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("deepl response status:", resp.Status)
		return 0, errors.New(resp.Status)
	}

	// Parse response
	var translateResp deeplTranslateResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&translateResp)
	if err != nil {
		return 0, err
	}

	log.Debugln("detected source language: ", translateResp.Translations[0].DetectedSourceLanguage)
	return io.Copy(w, strings.NewReader(translateResp.Translations[0].Text))
}
