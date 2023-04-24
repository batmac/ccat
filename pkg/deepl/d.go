package deepl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// https://www.deepl.com/fr/docs-api/translate-text/translate-text/

type Translation struct {
	DetectedSourceLanguage string `json:"detected_source_language"` //nolint:tagliatelle
	Text                   string `json:"text"`
}

type TranslateResponse struct {
	Translations []Translation `json:"translations"`
}

type Client struct {
	apiKey string
}

func NewDeepLClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

func (c *Client) Translate(text string, targetLang string, sourceLang string) (*TranslateResponse, error) {
	baseURL := "https://api-free.deepl.com/v2/translate"
	// if apikey doesnt end with ":fx", use paid API
	if c.apiKey[len(c.apiKey)-3:] != ":fx" {
		baseURL = "https://api.deepl.com/v2/translate"
	}

	// Check for text size limit
	if len(text) > 128*1024 {
		return nil, fmt.Errorf("text size exceeds 128 KiB limit")
	}

	// Create request body
	data := url.Values{}
	data.Set("text", text)
	data.Set("target_lang", targetLang)
	if sourceLang != "" {
		data.Set("source_lang", sourceLang)
	}

	// Create request
	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "DeepL-Auth-Key "+c.apiKey)

	// Send request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response
	var translateResp TranslateResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&translateResp)
	if err != nil {
		return nil, err
	}

	return &translateResp, nil
}
