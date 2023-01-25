package mutators_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
	"github.com/batmac/ccat/pkg/stringutils"
)

var testsFeed2y = []struct {
	name, decoded, encoded string
}{
	{
		"null", stringutils.DeleteSpaces(`
	feedType: atom
	feedVersion: "1.0"
	items: []
	`),
		`<?xml version="1.0" encoding="UTF-8"?>
	<feed xmlns="http://www.w3.org/2005/Atom">
	</feed>`,
	},
}

func Test_feed2Y(t *testing.T) {
	f := "feed2y"
	for _, tt := range testsFeed2y {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringutils.DeleteSpaces(mutators.Run(f, tt.encoded)); got != stringutils.DeleteSpaces(tt.decoded) {
				t.Errorf("%s = %v, want %v", f, got, tt.decoded)
			}
		})
	}
}
