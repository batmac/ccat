package term

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/utils"

	"braces.dev/errtrace"
	"golang.org/x/term"
)

func IsStdoutTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func IsStdinTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

func GetTerminalSize() (int, int, error) {
	if IsStdoutTerminal() {
		return errtrace.Wrap3(term.GetSize(int(os.Stdout.Fd())))
	}
	// fallback when piping to a file!
	return 80, 24, nil // VT100 terminal size
}

func ClearScreen() {
	if IsStdoutTerminal() {
		fmt.Print("\033[H\033[2J")
	}
}

func SupportedColors() uint {
	var colors uint
	if !IsStdoutTerminal() {
		log.Debugln("stdout is not a terminal, detecting colors anyway...")
	}

	switch {
	case IsITerm2():
		colors = 16_000_000
		log.Debugln("  supportedColors: iterm2 -> 16M colors detected")
	case strings.ToLower(os.Getenv("COLORTERM")) == "truecolor":
		colors = 16_000_000
		log.Debugln("  supportedColors: truecolor -> 16M colors detected")
	case os.Getenv("TERM") == "xterm-256color":
		colors = 256
		log.Debugln("  supportedColors: xterm-256color -> 256 colors detected")
	case utils.IsRunningInContainer():
		// if we're running in a container, we suppose the user has a sufficiently modern term.
		colors = 256
		log.Debugln("  supportedColors: container detected, setting to 256 colors")
	default:
		log.Debugf("  supportedColors: unknown term, $TERM==%s -> 8 colors detected\n", os.Getenv("TERM"))
		return 8
	}

	return colors
}

func ReadPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	b, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", errtrace.Wrap(err)
	}
	return strings.TrimSpace(string(b)), nil
}

func ReadLine(prompt string) (string, error) {
	fmt.Print(prompt)

	// read a line from stdin with echoing it to the terminal
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return "", errtrace.Wrap(err)
	}

	return strings.TrimSpace(line), nil
}
