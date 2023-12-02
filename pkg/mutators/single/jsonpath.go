package mutators

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"braces.dev/errtrace"
	"github.com/PaesslerAG/jsonpath"
	"github.com/batmac/ccat/pkg/log"
)

// https://goessner.net/articles/JsonPath/

func init() {
	singleRegister("jsonpath", jSONPath, withDescription("a jsonpath expression to apply (on $, with all ',' replaced by '|', all ':' replaced by '£')"),
		withConfigBuilder(stdConfigStringWithDefault("$")),
		withCategory("filter"),
	)
}

func jSONPath(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	buf, err := io.ReadAll(r) // NOT streamable
	if err != nil {
		return 0, errtrace.Wrap(err)
	}

	var v interface{}

	if err := json.Unmarshal(buf, &v); err != nil {
		return 0, errtrace.Wrap(err)
	}

	jp := config.(string)
	if jp[0] == '.' {
		jp = "$" + jp
	}
	// replace all "|" with ","
	jp = strings.ReplaceAll(jp, "|", ",")
	//
	jp = strings.ReplaceAll(jp, "£", ":")

	log.Debugf("final jsonpath: %s", jp)
	values, err := jsonpath.Get(jp, v)
	if err != nil {
		return 0, errtrace.Wrap(err)
	}

	// Marshal the result back to JSON
	buf, err = json.Marshal(values)
	if err != nil {
		return 0, errtrace.Wrap(err)
	}

	fmt.Fprintln(w, string(buf))

	return int64(len(buf)), nil
}
