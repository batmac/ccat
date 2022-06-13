package pipedcmd

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/batmac/ccat/log"
)

type PipedCmd struct {
	Cmd    *exec.Cmd
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser

	written int64
}

func New(cmdline ...string) (*PipedCmd, error) {
	log.Debugln(" pipedcmd: start new...")

	cmdline = strings.Split(strings.Join(cmdline, " "), " ")
	if cmdline[0] == "" {
		return nil, exec.ErrNotFound
	}

	ctx := context.TODO()

	fullCmd, err := exec.LookPath(cmdline[0])
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, fullCmd, cmdline[1:]...)
	/* 	cmd.Env = append(os.Environ(),
		"FOO=actual_value", // this value is used
	) */
	log.Debugf(" pipedcmd: created cmd %v", cmd)
	log.Debugln(" pipedcmd: defining stdstuff...")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	cmd.Stderr = os.Stderr
	stderr := os.Stderr
	/* stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	} */

	log.Debugln(" pipedcmd: returning new...")

	return &PipedCmd{
		Cmd:    cmd,
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
		//written: 0,
	}, nil
}

func (c *PipedCmd) Start(stdin io.Reader) error {
	log.Debugf(" pipedcmd: Start %v\n", c)

	err := c.Cmd.Start()
	if err != nil {
		return err
	}
	//print("read\n")

	go func() {
		log.Debugln(" pipedcmd: routine starting to copy... ")

		written, err := io.Copy(c.Stdin, stdin)
		if err != nil {
			log.Println(err)
		}
		c.written = written
		log.Debugf(" pipedcmd: routine written %d bytes to the cmd \"%s\"\n", written, c.Cmd)
		c.Stdin.Close()
		if err != nil {
			log.Println(err)
		}
		log.Debugln(" pipedcmd: routine end")
	}()
	return nil
}

// wait that cmd exits and copying ends.
func (c *PipedCmd) Wait() error {
	log.Debugf(" pipedcmd: Wait %v\n", c)
	return c.Cmd.Wait()
}

func (c *PipedCmd) String() string {
	return c.Cmd.String()
}
