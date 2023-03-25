package utils_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/utils"
)

func TestIsRunningIn(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		k8s := utils.IsRunningInK8s()
		docker := utils.IsRunningInDocker()
		podman := utils.IsRunningInPodman()
		_, got := utils.IsRunningInContainer(), utils.IsRunningInContainer()

		if got != (k8s || docker || podman) {
			t.Errorf(" got(%v) != (k8s(%v) || docker(%v))\n", got, k8s, docker)
		}
	})
}
