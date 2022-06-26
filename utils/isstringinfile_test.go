package utils_test

import (
	"os"
	"testing"

	"github.com/batmac/ccat/log"
	"github.com/batmac/ccat/utils"
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
		assertPanic(t, func() {
			_ = utils.IsStringInFile(" ", "not existing file")
		})
	})
	t.Run("panicplease", func(t *testing.T) {
		assertPanic(t, func() {
			_ = utils.IsStringInFile("", exe)
		})
	})
	t.Run("panic+notexistingfile", func(t *testing.T) {
		assertPanic(t, func() {
			_ = utils.IsStringInFile("", "not existing file")
		})
	})
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}
