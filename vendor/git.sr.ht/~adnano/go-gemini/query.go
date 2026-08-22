package gemini

import (
	"net/url"
	"strings"
)

// QueryEscape escapes a string for use in a Gemini URL query.
// It is like url.PathEscape except that it also replaces plus signs
// with their percent-encoded counterpart.
func QueryEscape(query string) string {
	return strings.ReplaceAll(url.PathEscape(query), "+", "%2B")
}

// QueryUnescape unescapes a Gemini URL query.
// It is identical to url.PathUnescape.
func QueryUnescape(query string) (string, error) {
	return url.PathUnescape(query)
}
