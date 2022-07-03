package main

import "testing"

func Test_setupStdout(t *testing.T) {
	t.Run("donotpanicplease", func(t *testing.T) {
		_ = setupStdout(false)
		_ = setupStdout(true)
	})
}
