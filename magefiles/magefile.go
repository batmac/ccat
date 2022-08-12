//nolint:deadcode
package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = BuildAndTest

var defaultBuildArgs = []string{"build"}

func ldFlags(goTags string) string {
	version, err := exec.Command("git", "describe", "--tags").Output()
	if err != nil {
		_ = mg.Fatal(1, err)
	}
	commit, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		_ = mg.Fatal(1, err)
	}

	return fmt.Sprintf("-s -w -X main.version=%s -X main.commit=%s -X main.date=%s -X main.builtBy=%s -X main.tags=%s",
		string(version),
		string(commit),
		time.Now().Format("2006-01-02@15:04:05-0700"),
		"Mage",
		goTags,
	)
}

func build(tags string) error {
	mg.Deps(InstallDeps)
	fmt.Println("Building...")
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir("cmd/ccat"); err != nil {
		return err
	}
	buildArgs := append(defaultBuildArgs, "-ldflags", ldFlags(tags))

	if err := sh.RunWithV(nil, "go", buildArgs...); err != nil {
		return err
	}
	if err := os.Rename("ccat", "../../ccat"); err != nil {
		return err
	}
	if err := os.Chdir(cwd); err != nil {
		return err
	}

	return nil
}

// tags: libcurl,crappy
func BuildDefault() error {
	return build("libcurl,crappy")
}

// tags: nohl,fileonly
func BuildMinimal() error {
	return build("nohl,fileonly")
}

// put ccat to $GOPATH/bin/ccat
func Install() error {
	mg.Deps(BuildAndTest)
	path := os.ExpandEnv("$GOPATH/bin/ccat")
	fmt.Printf("Installing to '%s'...\n", path)

	return sh.Copy(path, "ccat")
}

// go mod download
func InstallDeps() error {
	fmt.Println("Installing Deps...")
	return sh.RunV("go", "mod", "download")
}

// go mod verify
func VerifyDeps() error {
	mg.Deps(InstallDeps)
	fmt.Println("Verifying Deps...")
	return sh.Run("go", "mod", "verify")
}

func Clean() {
	fmt.Println("Cleaning...")
	_ = os.RemoveAll("ccat")
}

// go test ./...
func TestGo() error {
	mg.Deps(InstallDeps)
	fmt.Println("Testing Go...")
	return sh.RunV("go", "test", "./...")
}

// test_compression_e2e
func TestCompression() error {
	mg.Deps(InstallDeps)
	fmt.Println("Testing compression...")
	return sh.RunV("scripts/test_compression_e2e.sh", "testdata/compression/")
}

// test all
func Test() error {
	mg.SerialDeps(TestGo)
	mg.SerialDeps(TestCompression)
	return nil
}

// buildDefault,test
func BuildAndTest() error {
	mg.SerialDeps(BuildDefault)
	mg.SerialDeps(Test)
	return nil
}

func UpdateREADME() error {
	fmt.Println("Updating README.md...")

	data, err := os.ReadFile("README.header.md")
	if err != nil {
		return err
	}
	data = append(data, "\n```\n"...)
	cmd := exec.Command("./ccat", "--fullhelp")
	out, err := cmd.CombinedOutput()
	data = append(data, out...)
	if err != nil {
		return err
	}
	data = append(data, "```\n"...)

	return os.WriteFile("README.md", data, 0o600)
}
