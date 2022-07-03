package utils_test

import (
	"os"
	"testing"

	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/utils"
)

func TestIsStringInFile(t *testing.T) {
	exe, err := os.Executable() // we just need an existing file
	if err != nil {
		t.Fatal(err)
	}

	t.Run("donotpanicplease", func(t *testing.T) {
		got := utils.IsStringInFile(" ", exe)
		log.Debugf("%v\n", got)
	})
	t.Run("notexistingfile", func(t *testing.T) {
		if utils.IsStringInFile(" ", "not existing file") {
			t.Fail()
		}
	})
	t.Run("panicplease", func(t *testing.T) {
		assertPanic(t, func() {
			_ = utils.IsStringInFile("", exe)
		})
	})
}

func assertPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}
