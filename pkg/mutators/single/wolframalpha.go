package mutators

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"braces.dev/errtrace"
	"github.com/batmac/ccat/pkg/globalctx"
	"github.com/batmac/ccat/pkg/secretprovider"
)

func init() {
	singleRegister("wa", wolframalphashort,
		withDescription("query wolfram alpha Short Answers API (APPID in $WA_APPID)"),
		withCategory("external APIs"))

	singleRegister("waspoken", wolframalphaspoken,
		withDescription("query wolfram alpha Spoken API (APPID in $WA_APPID)"),
		withCategory("external APIs"))

	singleRegister("wasimple", wolframalphasimple,
		withDescription("query wolfram alpha Simple API (output is an image, APPID in $WA_APPID)"),
		withCategory("external APIs"),
		withExpectingBinary())

	// https://products.wolframalpha.com/llm-api/documentation
	singleRegister("wallm", wolframalphallm,
		withDescription("query wolfram alpha LLM API (APPID in $WA_APPID)"),
		withCategory("external APIs"))
}

func wolframalphashort(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	return errtrace.Wrap2(wolframalpha(w, r, "result", "i"))
}

func wolframalphaspoken(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	return errtrace.Wrap2(wolframalpha(w, r, "spoken", "i"))
}

func wolframalphasimple(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	return errtrace.Wrap2(wolframalpha(w, r, "simple", "i"))
}

func wolframalphallm(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	return errtrace.Wrap2(wolframalpha(w, r, "llm-api", "input"))
}

func wolframalpha(w io.WriteCloser, r io.ReadCloser, t string, queryField string) (int64, error) {
	baseURL := "https://api.wolframalpha.com/v1/" + t + "?"

	query, err := io.ReadAll(r) // NOT streamable
	if err != nil {
		return 0, errtrace.Wrap(err)
	}
	if len(query) == 0 {
		return 0, nil
	}

	appID, _ := secretprovider.GetSecret("wolfram", "WA_APPID")
	if appID == "" {
		return 0, errtrace.Wrap(fmt.Errorf("no appid found in $WA_APPID"))
	}
	// build query url values
	q := url.Values{
		queryField: {string(query)},
		"appid":    {appID},
	}

	res, err := http.Get(baseURL + q.Encode())
	if err != nil {
		return 0, errtrace.Wrap(err)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, errtrace.Wrap(err)
	}
	res.Body.Close()
	// fmt.Println(string(data))

	expectingBinary := globalctx.Get("expectingBinary")
	// ensure data last char is a newline
	if (expectingBinary == nil || !expectingBinary.(bool)) && data[len(data)-1] != '\n' {
		data = append(data, '\n')
	}
	return errtrace.Wrap2(io.Copy(w, bytes.NewReader(data)))
}
