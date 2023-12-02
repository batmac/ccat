package plugins

import (
	"braces.dev/errtrace"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// ccat -m yaegi:example/plugins/deepl.go:plugins.Deepl echo://"bonjour"

type deeplTranslation struct {
	DetectedSourceLanguage string `json:"detected_source_language"` //nolint:tagliatelle
	Text                   string `json:"text"`
}

type deeplTranslateResponse struct {
	Translations []deeplTranslation `json:"translations"`
}

func Deepl(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	targetLanguage := os.Getenv("TARGET_LANGUAGE")
	if targetLanguage == "" {
		targetLanguage = "en"
	}

	text, err := io.ReadAll(r) // NOT streamable
	if err != nil {
		return 0, errtrace.Wrap(err)
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
		return 0, errtrace.Wrap(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "DeepL-Auth-Key "+key)

	// Send request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, errtrace.Wrap(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("deepl response status:", resp.Status)
		return 0, errtrace.Wrap(errors.New(resp.Status))
	}

	// Parse response
	var translateResp deeplTranslateResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&translateResp)
	if err != nil {
		return 0, errtrace.Wrap(err)
	}

	log.Println("detected source language: ", translateResp.Translations[0].DetectedSourceLanguage)
	return errtrace.Wrap2(io.Copy(w, strings.NewReader(translateResp.Translations[0].Text)))
}
