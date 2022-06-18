package term

import (
	"encoding/base64"
	"fmt"

	"github.com/batmac/ccat/log"
)

func Osc52(d []byte) {
	log.Debugf("writing to clipboard via osc52\n")
	fmt.Print("\033]52;c;")
	fmt.Print(base64.StdEncoding.EncodeToString(d))
	fmt.Print("\a\n")
}
