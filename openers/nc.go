//go:build !fileonly
// +build !fileonly

package openers

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/batmac/ccat/log"
	"github.com/batmac/ccat/utils"
)

var (
	tcpOpenerName        = "tcp"
	tcpOpenerDescription = "get data from listening on tcp://[HOST]:<PORT>"
)

type tcpOpener struct {
	name, description string
}

func init() {
	register(&tcpOpener{
		name:        tcpOpenerName,
		description: tcpOpenerDescription,
	})
}

func (f tcpOpener) Name() string {
	return f.name
}

func (f tcpOpener) Description() string {
	return f.description
}

func (f tcpOpener) Open(s string, _ bool) (io.ReadCloser, error) {
	l, err := net.Listen("tcp", utils.RemoveScheme(s))
	if err != nil {
		return nil, fmt.Errorf("Error listening: %v", err)
	}
	log.Debugln("Listening on " + utils.RemoveScheme(s))

	conn, err := l.Accept()
	if err != nil {
		return nil, fmt.Errorf("Error accepting: %v", err)
	}
	log.Debugln("Accepted, closing the listening socket...")
	l.Close()

	return conn, nil
}

func (f tcpOpener) Evaluate(s string) float32 {
	// log.Debugf("Evaluating %s...\n", s)
	if strings.HasPrefix(s, "tcp://") {
		return 0.9
	}
	return 0
}
