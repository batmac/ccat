package term

import (
	"os"

	"github.com/batmac/ccat/pkg/log"
)

func IsSsh() bool {
	if len(os.Getenv("SSH_TTY")) > 0 || len(os.Getenv("SSH_CONNECTION")) > 0 || len(os.Getenv("SSH_CLIENT")) > 0 {
		log.Debugf("ssh detected\n")
		return true
	}
	return false
}
