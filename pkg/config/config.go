package config

import (
	"os"
	"path/filepath"

	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/utils"
	"gopkg.in/yaml.v2"
)

const (
	appName               = "ccat"
	defaultConfigFilename = appName + ".conf"

	currentConfigVersion = "v1"
)

func Exists(path string) bool {
	return utils.PathExists(path)
}

func Path() string {
	if os.Getenv("CCATRC") != "" {
		log.Debugln("Using config path from $CCATRC")
		return os.Getenv("CCATRC")
	}

	if Exists("~/.ccatrc") {
		log.Debugln("Using config path ~/.ccatrc")
		return "~/.ccatrc"
	}

	if os.Getenv("XDG_CONFIG_HOME") != "" {
		path := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), appName, defaultConfigFilename)
		if Exists(path) {
			log.Debugln("Using config path from $XDG_CONFIG_HOME")
			return path
		}
	}

	var path string
	switch os.Getenv("GOOS") {
	case "windows":
		path = filepath.Join(os.Getenv("APPDATA"), appName, defaultConfigFilename)
	case "darwin":
		path = filepath.FromSlash("~/Library/Application Support/" + appName + "/" + defaultConfigFilename)
	default: // includes unix platforms
		path = filepath.FromSlash("~/.config/" + appName + "/" + defaultConfigFilename)
	}
	if Exists(path) {
		log.Debugln("Using config path " + path)
		return path
	}

	log.Debugln("No config file found")
	return ""
}

type Config struct {
	Aliases map[string]string `yaml:"aliases"`
	Version string            `yaml:"version"`
}

func New() *Config {
	return &Config{
		Version: currentConfigVersion,
		Aliases: make(map[string]string),
	}
}

func Load() *Config {
	c := New()
	expandedPath := utils.ExpandPath(Path())
	if expandedPath == "" {
		log.Debugf("config path not determined, not loading")
		return c
	}
	log.Debugln("Loading config from " + expandedPath)
	raw, err := os.ReadFile(expandedPath)
	if err != nil {
		log.Printf("error loading %s: %s\n", expandedPath, err)
		return c
	}
	if err := yaml.Unmarshal(raw, &c); err != nil {
		log.Printf("error parsing %s: %s\n", expandedPath, err)
		return c
	}
	if c.Version != currentConfigVersion {
		log.Printf("config version %s not supported, ignoring config\n", c.Version)
		return New()
	}
	log.Debugf("Loaded config: %+v", c)
	return c
}

func Save(c *Config) error {
	if c == nil {
		log.Debugf("no config to save")
		return nil
	}
	expandedPath := utils.ExpandPath(Path())
	if expandedPath == "" {
		log.Debugf("config path not determined, not saving")
		return nil
	}
	log.Debugln("Saving config to " + expandedPath)
	raw, err := yaml.Marshal(c)
	if err != nil {
		log.Printf("error marshalling config: %s\n", err)
		return err
	}
	//#nosec G306
	if err := os.WriteFile(expandedPath, raw, 0o644); err != nil {
		log.Printf("error writing %s: %s\n", expandedPath, err)
		return err
	}
	return nil
}
