package utils

import (
	"log"
	"net/url"
)

func RemoveScheme(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		log.Fatal(err)
	}
	return u.Host
}
