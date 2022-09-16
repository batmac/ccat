//nolint:deadcode // obvious for mage
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/batmac/ccat/pkg/term"
	"github.com/google/renameio/maybe"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = BuildDefaultAndTest

var (
	defaultBuildArgs = []string{"build"}
	binaryName       = "ccat"
)

func init() {
	if term.IsStdoutTerminal() {
		os.Setenv("MAGEFILE_ENABLE_COLOR", "1")
	}
	if runtime.GOOS == "windows" {
		binaryName = "ccat.exe"
	}
}

func ldFlags(goTags string) string {
	version, err := sh.Output("git", "describe", "--tags")
	if err != nil {
		_ = mg.Fatal(1, err)
	}
	commit, err := sh.Output("git", "rev-parse", "HEAD")
	if err != nil {
		_ = mg.Fatal(1, err)
	}

	return fmt.Sprintf("-s -w -X main.version=%s -X main.commit=%s -X main.date=%s -X main.builtBy=%s -X main.tags=%s",
		version,
		commit,
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
	if err := os.Chdir(filepath.FromSlash("cmd/ccat")); err != nil {
		return err
	}
	buildArgs := defaultBuildArgs
	buildArgs = append(buildArgs, "-ldflags", ldFlags(tags), "-tags", tags)

	if err := sh.RunWithV(nil, "go", buildArgs...); err != nil {
		return err
	}

	if err := os.Rename(binaryName, filepath.Join(cwd, binaryName)); err != nil {
		return err
	}
	if err := os.Chdir(cwd); err != nil {
		return err
	}

	stepOKPrintln("Building OK")
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
	path := os.ExpandEnv(filepath.Join("$GOPATH", "bin", binaryName))
	stepPrintf("Installing... (%s)\n", path)

	data, err := os.ReadFile(binaryName)
	if err != nil {
		// detect this
		fmt.Println("Have you build first?")
		return err
	}
	if err := maybe.WriteFile(path, data, 0o750); err != nil {
		return err
	}
	stepOKPrintln("Installing OK")
	return nil
}

// go mod download
func InstallDeps() error {
	stepPrintln("Installing Deps...")
	if err := sh.RunV("go", "mod", "download"); err != nil {
		return err
	}
	stepOKPrintln("Installing Deps OK")
	return nil
}

// go mod verify
func VerifyDeps() error {
	mg.Deps(InstallDeps)
	stepPrintln("Verifying Deps...")
	if err := sh.Run("go", "mod", "verify"); err != nil {
		return err
	}
	stepOKPrintln("Verifying Deps OK")
	return nil
}

func Clean() error {
	stepPrintln("Cleaning...")
	if err := sh.Rm(binaryName); err != nil {
		return err
	}
	stepOKPrintln("Cleaning OK")
	return nil
}

// go test ./...
func TestGo() error {
	mg.Deps(InstallDeps)
	stepPrintln("Testing Go...")
	r, err := sh.Output("go", "test", "./...")
	if mg.Verbose() {
		fmt.Println(r, "\n ")
	}
	if err != nil {
		return err
	}
	stepOKPrintln("Testing Go OK")
	return nil
}

// test_compression_e2e
/* func TestCompression() error {
	stepPrintln("Testing compression...")
	return sh.RunV("scripts/test_compression_e2e.sh", "testdata/compression/")
} */

// test all
func Test() error {
	mg.SerialDeps(TestGo)
	mg.SerialDeps(TestCompressionGo)
	return nil
}

// buildDefault,test
func BuildDefaultAndTest() error {
	mg.SerialDeps(BuildDefault)
	mg.SerialDeps(Test)
	return nil
}

func UpdateREADME() error {
	stepPrintln("Updating README.md...")

	data, err := os.ReadFile("README.header.md")
	if err != nil {
		return err
	}
	data = append(data, "\n```\n"...)
	out, err := sh.Output("./"+binaryName, "--fullhelp")

	data = append(data, out...)
	if err != nil {
		return err
	}
	data = append(data, "```\n"...)

	if err := os.WriteFile("README.md", data, 0o600); err != nil {
		return err
	}
	stepOKPrintln("Updating README.md OK")
	return nil
}

func stepPrintln(a ...any) {
	fmt.Println(append([]any{"🚧"}, a...)...)
}

func stepPrintf(format string, a ...any) {
	fmt.Printf("🚧 "+format, a...)
}

func stepOKPrintln(a ...any) {
	prefix := ""
	if !mg.Verbose() {
		prefix = "\x1bM" // go on the previous line
	}
	fmt.Println(append([]any{prefix + "✅"}, a...)...)
}
