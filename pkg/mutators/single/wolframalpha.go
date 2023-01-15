package mutators

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/batmac/ccat/pkg/globalctx"
)

func init() {
	singleRegister("wa", wolframalphashort, withDescription("query wolfram alpha Short Answers API (APPID in $WA_APPID)"))
	singleRegister("waspoken", wolframalphaspoken, withDescription("query wolfram alpha Spoken API (APPID in $WA_APPID)"))
	singleRegister("wasimple", wolframalphasimple, withDescription("query wolfram alpha Simple API (output is an image, APPID in $WA_APPID)"),
		withExpectingBinary(true))
}

func wolframalphashort(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	return wolframalpha(w, r, "result")
}

func wolframalphaspoken(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	return wolframalpha(w, r, "spoken")
}

func wolframalphasimple(w io.WriteCloser, r io.ReadCloser, _ any) (int64, error) {
	return wolframalpha(w, r, "simple")
}

func wolframalpha(w io.WriteCloser, r io.ReadCloser, t string) (int64, error) {
	baseURL := "https://api.wolframalpha.com/v1/" + t + "?"

	query, err := io.ReadAll(r) // NOT streamable
	if err != nil {
		return 0, err
	}
	if len(query) == 0 {
		return 0, nil
	}

	appID := os.Getenv("WA_APPID")
	if appID == "" {
		return 0, fmt.Errorf("no appid found in $WA_APPID")
	}
	// build query url values
	q := url.Values{
		"i":     {string(query)},
		"appid": {appID},
	}

	res, err := http.Get(baseURL + q.Encode())
	if err != nil {
		return 0, err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	res.Body.Close()
	// fmt.Println(string(data))

	expectingBinary := globalctx.Get("expectingBinary")
	// ensure data last char is a newline
	if (expectingBinary == nil || !expectingBinary.(bool)) && data[len(data)-1] != '\n' {
		data = append(data, '\n')
	}
	return io.Copy(w, bytes.NewReader(data))
}
