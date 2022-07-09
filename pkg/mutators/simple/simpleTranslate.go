package mutators

import (
	"encoding/json"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/batmac/ccat/pkg/log"
)

const postUrl = "https://translation.googleapis.com/language/translate/v2"

func init() {
	simpleRegister("translate", translate,
		withDescription("translate to $TARGET_LANGUAGE with google translate (need a valid key in $GOOGLE_API_KEY)"),
	)
}

type Response struct {
	Data struct {
		Translations []struct {
			TranslatedText         string
			DetectedSourceLanguage string
		}
	}
}

func translate(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	msg, err := ioutil.ReadAll(r) // NOT streamable
	if err != nil {
		return 0, err
	}

	if len(msg) == 0 {
		return 0, nil
	}

	result := strings.Builder{}

	key := os.Getenv("GOOGLE_API_KEY")
	if key == "" {
		log.Fatal("no key found in $GOOGLE_API_KEY")
	}
	targetLanguage := "en"
	if t := os.Getenv("TARGET_LANGUAGE"); t != "" {
		targetLanguage = t
	}

	v := url.Values{}
	v.Set("key", key)
	v.Set("q", string(msg))
	v.Set("target", targetLanguage)

	res, err := http.PostForm(postUrl, v)
	if err != nil {
		return 0, err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	res.Body.Close()
	// fmt.Println(string(data))

	var d Response
	if err := json.Unmarshal(data, &d); err != nil {
		return 0, err
	}
	if len(d.Data.Translations) == 0 {
		log.Printf("didn't get any translation, got %s\n", string(data))
		return 0, nil
	}
	log.Debugf("Found a translation from language %s to %s\n", d.Data.Translations[0].DetectedSourceLanguage, targetLanguage)
	result.WriteString(html.UnescapeString(d.Data.Translations[0].TranslatedText))
	for _, txt := range d.Data.Translations[1:] {
		log.Debugf("Found an extra translation from language %s: %s\n", txt.DetectedSourceLanguage, txt.TranslatedText)
	}

	return io.Copy(w, strings.NewReader(result.String()))
}
