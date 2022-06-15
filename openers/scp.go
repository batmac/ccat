//go:build crappy && !fileonly
// +build crappy,!fileonly

package openers

import (
	"github.com/batmac/ccat/log"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var localShellScpOpenerName = "ShellScp"
var localShellScpOpenerDescription = "get scp:// via local scp\n"

type localShellScpOpener struct {
	name, description string
}

func init() {
	_ = register(&localShellScpOpener{
		name:        localShellScpOpenerName,
		description: localShellScpOpenerDescription,
	})
}

func (f localShellScpOpener) Name() string {
	return f.name
}
func (f localShellScpOpener) Description() string {
	return f.description
}
func (f *localShellScpOpener) Open(s string, _ bool) (io.ReadCloser, error) {

	log.Debugln(" localShellScp started")

	arr := strings.SplitN(s, "scp://", 2)
	path := strings.Split(arr[1], " ")
	tmpfile, err := os.CreateTemp("", "ccat_tempfile_")
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf(" localShellScp temp file is %s\n", tmpfile.Name())
	path = append(path, tmpfile.Name())

	cmd := exec.Command("scp", path...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	so, err := ioutil.ReadAll(stdout)
	se, err := ioutil.ReadAll(stderr)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugln(string(so), string(se))
	log.Debugf(" localShellScp ended\n")
	re, err := ioutil.ReadAll(tmpfile)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpfile.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = os.Remove(tmpfile.Name())
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	return ioutil.NopCloser(strings.NewReader(string(re))), nil
}

func (f localShellScpOpener) Evaluate(s string) float32 {
	arr := strings.SplitN(s, "scp://", 2)
	before := arr[0]
	if before == "" {
		return 0.5
	}
	return 0
}
