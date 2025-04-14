package selfupdate

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/batmac/ccat/pkg/log"
	"github.com/creativeprojects/go-selfupdate"
)

type Mode int8

const (
	ModeCheckOnly Mode = iota
	ModeUpdate
	ModeForce
)

// build tags for the github releases
var (
	githubTags        = "plugins,keystore"
	tagsAreCompatible = false
)

func Do(version, tags string, mode Mode) {
	log.Debugf("Trying to self-update %v...\n", version)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	selfupdate.SetLogger(log.Debug)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("self-update failed: ", r)
			os.Exit(98)
		}
	}()

	updater, _ := selfupdate.NewUpdater(selfupdate.Config{Validator: &selfupdate.ChecksumValidator{UniqueFilename: "checksums.txt"}})

	latest, found, err := updater.DetectLatest(ctx, selfupdate.NewRepositorySlug("batmac", "ccat"))
	if err != nil {
		panic(fmt.Errorf("error occurred while detecting version: %w", err))
	}
	if !found {
		panic(fmt.Errorf("latest version for %s/%s could not be found from github repository", runtime.GOOS, runtime.GOARCH))
	}

	cleanedVersion := CleanVersion(version)
	log.Debugf("cleaned version is '%v'\n", cleanedVersion)
	if latest.LessOrEqual(cleanedVersion) {
		fmt.Printf("Current version (%s) is the latest\n", version)
		if mode == ModeForce {
			fmt.Println("But you forced me to update, so I will!")
		} else {
			return
		}
	}

	fmt.Printf("Version %v is available\n", latest.Version())

	if tags != githubTags {
		fmt.Printf("Warning: your current binary is built with tags '%s', GitHub releases are built with '%s'.\n", tags, githubTags)
		tagsAreCompatible = false
	} else {
		tagsAreCompatible = true
	}

	if mode == ModeCheckOnly {
		return
	}

	if !tagsAreCompatible {
		fmt.Println("I'm about to update your binary with the last one available from GitHub, " +
			"which doesn't have the same build tags, this may not be what you want, " +
			"hit Ctrl-C to abort or Enter to continue.")
		input := make(chan string)

		go func() {
			_, _ = fmt.Scanln()
			close(input)
		}()
		select {
		case <-input:
		case <-time.After(300 * time.Second):
			fmt.Println("timed out, aborting.")
			//nolint:gocritic // exitAfterDefer
			os.Exit(97)
		}
	}

	exe, err := os.Executable()
	if err != nil {
		panic(errors.New("could not locate executable path"))
	}

	if err := updater.UpdateTo(ctx, latest, exe); err != nil {
		panic(fmt.Errorf("error occurred while updating binary: %w", err))
	}
	fmt.Printf("Successfully updated to version %s\n", latest.Version())
}

func CleanVersion(v string) string {
	// only keep the semver string on patch/git version
	v, _, _ = strings.Cut(v, "-")
	s := []byte(v)
	j := 0
	for _, b := range s {
		if ('0' <= b && b <= '9') ||
			b == '.' {
			s[j] = b
			j++
		}
	}
	c := string(s[:j])
	if c == "" {
		return "0"
	}
	return c
}
