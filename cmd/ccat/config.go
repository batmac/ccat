package main

import (
	"os"
	"path/filepath"

	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/utils"
)

const (
	appName               = "ccat"
	defaultConfigFilename = appName + ".conf"
)

func ConfigReadAll(path string) []byte {
	c, err := os.ReadFile(utils.ExpandPath(path))
	if err != nil {
		log.Printf("error loading %s: %s\n", path, err)
	}
	return c
}

func ConfigExists(path string) bool {
	_, err := os.Stat(utils.ExpandPath(path))
	return err == nil
}

func ConfigContent() []byte {
	if os.Getenv("CCATRC") != "" {
		log.Debugln("Using config from $CCATRC")
		return ConfigReadAll(os.Getenv("CCATRC"))
	}

	if ConfigExists("~/.ccatrc") {
		log.Debugln("Using config from ~/.ccatrc")
		return ConfigReadAll("~/.ccatrc")
	}

	if os.Getenv("XDG_CONFIG_HOME") != "" {
		path := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), appName, defaultConfigFilename)
		if ConfigExists(path) {
			log.Debugln("Using config from $XDG_CONFIG_HOME")
			return ConfigReadAll(path)
		}
	}

	var path string
	switch os.Getenv("GOOS") {
	case "windows":
		path = filepath.Join(os.Getenv("APPDATA"), appName, defaultConfigFilename)
	case "darwin":
		path = filepath.Join("~/Library/Application Support/", appName, defaultConfigFilename)
	default: // includes unix platforms
		path = filepath.Join("~/.config", appName, defaultConfigFilename)
	}
	if ConfigExists(path) {
		log.Debugln("Using config from " + path)
		return ConfigReadAll(path)
	}

	log.Debugln("No config file found")
	return nil
}
