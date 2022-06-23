package utils_test

import (
	"testing"

	"github.com/batmac/ccat/utils"
)

func TestIsRunningIn(t *testing.T) {
	t.Run("donotpanicplease", func(t *testing.T) {
		got := utils.IsRunningInContainer()
		k8s := utils.IsRunningInK8s()
		docker := utils.IsRunningInDocker()

		if got != (k8s || docker) {
			t.Errorf(" got(%v) != (k8s(%v) || docker(%v))\n", got, k8s, docker)
		}
	})
}
