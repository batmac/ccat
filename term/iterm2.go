package term

import (
	"ccat/log"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
)

func IsITerm2() bool {
	// LC_TERMINAL = iTerm2
	// TERM_PROGRAM = iTerm.app
	if os.Getenv("TERM_PROGRAM") == "iTerm.app" {
		return IsStdoutTerminal()
	}
	return false
}

func printITerm2ArtFromURL(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	if err := printITerm2Art(res.Body); err != nil {
		log.Println(err)
	}
}

func printITerm2Art(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	fmt.Print("\033]1337;")
	fmt.Print("File=inline=1;")
	fmt.Print("preserveAspectRatio=1;")
	fmt.Print(":")
	fmt.Print(base64.StdEncoding.EncodeToString(data))
	fmt.Print("\a\n")

	return nil
}
