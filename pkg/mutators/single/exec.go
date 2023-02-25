package mutators

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/batmac/ccat/pkg/log"
)

func init() {
	singleRegister("x", pipedcmd, withDescription("execute command (e.g. 'x:head -n 10')"),
		withConfigBuilder(stdConfigStringWithDefault("")))
}

func pipedcmd(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	command := config.(string)

	if command == "" {
		return 0, fmt.Errorf("no command specified")
	}

	split := strings.Split(command, " ")

	log.Debugf("executing command '%v'\n", split)
	cmd := exec.Command(split[0], split[1:]...) //nolint:gosec // that's the point
	cmd.Stdin = r
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return 0, err
}
