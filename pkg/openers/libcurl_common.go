//go:build !fileonly && ((cgo && libcurl) || (darwin && libcurl_purego))

package openers

import (
	"regexp"
	"strings"

	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/stringutils"
)

// shared by the cgo (libcurl) and the purego (libcurl_purego) implementations

const curlOpenerName = "curl"

func curlEvaluate(s string, protocols []string) float32 {
	// https://everything.curl.dev/protocols/curl
	// The latest curl (as of this writing) supports these protocols:
	// DICT, FILE, FTP, FTPS, GOPHER, GOPHERS, HTTP, HTTPS, IMAP, IMAPS, LDAP, LDAPS,
	// MQTT, POP3, POP3S, RTMP, RTSP, SCP, SFTP, SMB, SMBS, SMTP, SMTPS, TELNET, TFTP
	arr := strings.SplitN(s, "://", 2)
	before := arr[0]
	// log.Printf("before=%s found=%v s=%v", before, found, s)
	if stringutils.IsStringInSlice(before, protocols) {
		return 0.1
	}
	if before == "curlhttp" || before == "curlhttps" {
		return 1.0
	}
	return 0
}

func tryTransformURL(s string) string {
	// ease life by checking urls

	r := regexp.MustCompile(`^https://github.com/(.+)/blob(/.+)$`)
	matches := r.FindStringSubmatch(s)

	if len(matches) > 0 {
		url := "https://raw.githubusercontent.com/" + matches[1] + matches[2]
		log.Debugf("%s looks like a github url, transforming it to get the raw version: %s", s, url)
		return url
	}

	if strings.HasPrefix(s, "curlhttp://") || strings.HasPrefix(s, "curlhttps://") {
		return strings.TrimPrefix(s, "curl")
	}
	return s
}
