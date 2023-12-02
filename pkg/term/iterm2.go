package term

import (
	"braces.dev/errtrace"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func IsITerm2() bool {
	// LC_TERMINAL = iTerm2
	// TERM_PROGRAM = iTerm.app
	if os.Getenv("TERM_PROGRAM") == "iTerm.app" && os.Getenv("IGNORE_ITERM2") == "" {
		return IsStdoutTerminal()
	}
	return false
}

/* func printITerm2ArtFromURL(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	if err := printITerm2Art(res.Body); err != nil {
		log.Println(err)
	}
} */

func PrintITerm2Art(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return errtrace.Wrap(err)
	}

	fmt.Print("\033]1337;")
	fmt.Print("File=inline=1;")
	fmt.Print("preserveAspectRatio=1;")
	fmt.Print(":")
	fmt.Print(base64.StdEncoding.EncodeToString(data))
	fmt.Print("\a\n")

	return nil
}
