package utils

import (
	"os"

	"github.com/mitchellh/go-homedir"
)

// from https://github.com/charmbracelet/glow/blob/master/utils/utils.go (MIT)

// Expands tilde and all environment variables from the given path.
func ExpandPath(path string) string {
	s, err := homedir.Expand(path)
	if err == nil {
		return os.ExpandEnv(s)
	}
	return os.ExpandEnv(path)
}

func HomeDir() string {
	p, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	return p
}
