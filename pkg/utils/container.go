package utils

import (
	"github.com/batmac/ccat/pkg/log"
	"github.com/kofalt/go-memoize"
	"github.com/patrickmn/go-cache"
)

var globalCache = memoize.NewMemoizer(cache.NoExpiration, cache.NoExpiration)

func IsRunningInContainer() bool {
	log.Debugf("IsRunningInContainer?")
	return IsRunningInDocker() || IsRunningInPodman() || IsRunningInK8s()
}

func IsRunningInDocker() bool {
	result, _, _ := globalCache.Memoize("isRunningInDocker", isRunningInDocker)
	// fmt.Printf("isRunningInDocker cached: %v\n", cached)
	return result.(bool)
}

func IsRunningInPodman() bool {
	result, _, _ := globalCache.Memoize("isRunningInPodman", isRunningInPodman)
	// fmt.Printf("isRunningInPodman cached: %v\n", cached)
	return result.(bool)
}

func IsRunningInK8s() bool {
	result, _, _ := globalCache.Memoize("isRunningInK8s", isRunningInK8s)
	// fmt.Printf("isRunningInK8s cached: %v\n", cached)
	return result.(bool)
}

func isRunningInDocker() (any, error) {
	if PathExists("/.dockerenv") {
		log.Debugf("docker detected.")
		return true, nil
	}
	return false, nil
}

func isRunningInPodman() (any, error) {
	if PathExists("/run/.containerenv") {
		log.Debugf("podman detected.")
		return true, nil
	}
	return false, nil
}

func isRunningInK8s() (any, error) {
	// this does not work if automountServiceAccountToken: false
	if PathExists("/run/secrets/kubernetes.io/") {
		log.Debugf("k8s (at least one secret is mounted) detected.")
		return true, nil
	}
	if IsStringInFile("Kubernetes-managed hosts file", "/etc/hosts") {
		log.Debugf("k8s (managed /etc/hosts) detected.")
		return true, nil
	}
	return false, nil
}
