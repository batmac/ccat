package utils

import (
	"testing"

	"github.com/batmac/ccat/log"
)

func TestIsStringInFile(t *testing.T) {
	t.Run("donotpanicplease", func(t *testing.T) {
		got := IsStringInFile("", "not existing file")
		log.Debugf("%v\n", got)
	})
}
