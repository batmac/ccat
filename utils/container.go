package utils

import (
	"os"

	"github.com/batmac/ccat/log"
)

func IsRunningInContainer() bool {
	log.Debugf("IsRunningInContainer?")
	return IsRunningInDocker() || IsRunningInK8s()
}

func IsRunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		log.Debugf("docker detected.")
		return true
	}
	return false
}

func IsRunningInK8s() bool {
	// these does not work if automountServiceAccountToken: false
	if _, err := os.Stat("/run/secrets/kubernetes.io/"); err == nil {
		log.Debugf("k8s (at least one secret is mounted) detected.")
		return true
	}
	if IsStringInFile("Kubernetes-managed hosts file", "/etc/hosts") {
		log.Debugf("k8s (managed /etc/hosts) detected.")
		return true
	}
	return false
}
