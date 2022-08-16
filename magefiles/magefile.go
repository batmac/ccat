//nolint:deadcode
package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/google/renameio/maybe"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = BuildAndTest

var (
	defaultBuildArgs = []string{"build"}
	binaryName       = "ccat"
)

func init() {
	if runtime.GOOS == "windows" {
		binaryName = "ccat.exe"
	}
}

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
	stepPrintln("Building...")
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir("cmd/ccat"); err != nil {
		return err
	}
	buildArgs := append(defaultBuildArgs, "-ldflags", ldFlags(tags), "-tags", tags)

	if err := sh.RunWithV(nil, "go", buildArgs...); err != nil {
		return err
	}

	if err := os.Rename(binaryName, "../../"+binaryName); err != nil {
		return err
	}
	if err := os.Chdir(cwd); err != nil {
		return err
	}

	return nil
}

func BuildDefault() error {
	return build("")
}

// tags: nohl,fileonly
func BuildMinimal() error {
	return build("nohl,fileonly")
}

// tags: libcurl,crappy
func BuildFull() error {
	return build("libcurl,crappy")
}

// put ccat to $GOPATH/bin/ccat
func Install() error {
	mg.Deps(BuildAndTest)

	path := os.ExpandEnv("$GOPATH/bin/" + binaryName)
	stepPrintf("Installing to '%s'...\n", path)

	data, err := os.ReadFile(binaryName)
	if err != nil {
		return err
	}
	return maybe.WriteFile(path, data, 0o750)
}

// go mod download
func InstallDeps() error {
	stepPrintln("Installing Deps...")
	return sh.RunV("go", "mod", "download")
}

// go mod verify
func VerifyDeps() error {
	mg.Deps(InstallDeps)
	stepPrintln("Verifying Deps...")
	return sh.Run("go", "mod", "verify")
}

func Clean() {
	stepPrintln("Cleaning...")
	_ = os.RemoveAll(binaryName)
}

// go test ./...
func TestGo() error {
	mg.Deps(InstallDeps)
	stepPrintln("Testing Go...")
	return sh.RunV("go", "test", "./...")
}

// test_compression_e2e
func TestCompression() error {
	mg.Deps(InstallDeps)
	stepPrintln("Testing compression...")
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
	stepPrintln("Done.")
	return nil
}

func UpdateREADME() error {
	stepPrintln("Updating README.md...")

	data, err := os.ReadFile("README.header.md")
	if err != nil {
		return err
	}
	data = append(data, "\n```\n"...)
	cmd := exec.Command("./"+binaryName, "--fullhelp")
	out, err := cmd.CombinedOutput()
	data = append(data, out...)
	if err != nil {
		return err
	}
	data = append(data, "```\n"...)

	return os.WriteFile("README.md", data, 0o600)
}

func stepPrintln(a ...any) {
	fmt.Println(append([]any{"🚧"}, a...)...)
}

func stepPrintf(format string, a ...any) {
	fmt.Printf("🚧 "+format, a...)
}
