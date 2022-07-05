package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/batmac/ccat/pkg/log"
	"github.com/creativeprojects/go-selfupdate"
)

func update(version string, checkOnly bool) {
	log.Debugf("Trying to self-update %v...\n", version)

	selfupdate.SetLogger(log.Debug)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("self-update failed: ", r)
			os.Exit(98)
		}
	}()

	updater, _ := selfupdate.NewUpdater(selfupdate.Config{Validator: &selfupdate.ChecksumValidator{UniqueFilename: "checksums.txt"}})

	latest, found, err := updater.DetectLatest("batmac/ccat")
	if err != nil {
		panic(fmt.Errorf("error occurred while detecting version: %v", err))
	}
	if !found {
		panic(fmt.Errorf("latest version for %s/%s could not be found from github repository", runtime.GOOS, runtime.GOARCH))
	}

	cleanedVersion := cleanVersion(version)
	log.Debugf("cleaned version is %v\n", cleanedVersion)
	if latest.LessOrEqual(cleanedVersion) {
		log.Printf("Current version (%s) is the latest", version)
		return
	}

	fmt.Printf("Update to version %v is available\n", latest.Version())

	if checkOnly {
		return
	}

	exe, err := os.Executable()
	if err != nil {
		panic(errors.New("could not locate executable path"))
	}

	if err := updater.UpdateTo(latest, exe); err != nil {
		panic(fmt.Errorf("error occurred while updating binary: %v", err))
	}
	fmt.Printf("Successfully updated to version %s\n", latest.Version())
}

func cleanVersion(v string) string {
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
	return string(s[:j])
}
