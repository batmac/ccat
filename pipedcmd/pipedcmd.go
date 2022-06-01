package pipedcmd

import (
	"context"
	"io"
	"log"
	"os/exec"
	"strings"
)

type PipedCmd struct {
	Cmd    *exec.Cmd
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser

	written int64
}

func New(cmdline ...string) (*PipedCmd, error) {
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

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	/* stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	} */

	return &PipedCmd{
		Cmd:    cmd,
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: nil,
		//written: 0,
	}, nil
}

func (c *PipedCmd) Start(stdin io.Reader) error {
	err := c.Cmd.Start()
	if err != nil {
		return err
	}
	//print("read\n")

	go func() {
		written, err := io.Copy(c.Stdin, stdin)
		if err != nil {
			log.Println(err)
		}
		c.written = written
		//log.Printf("written %d bytes to the cmd \"%s\"\n", written, c.Cmd)
		c.Stdin.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	return nil
}

// wait that cmd exits and copying ends.
func (c *PipedCmd) Wait() error {
	//log.Printf("Wait %v\n", c)
	return c.Cmd.Wait()
}

func (c *PipedCmd) String() string {
	return c.Cmd.String()
}
