package mutators

import (
	"io"

	"github.com/batmac/ccat/pkg/log"

	"github.com/docker/go-units"
)

func init() {
	singleRegister("limit", limit, withDescription("a simple limiting fifo (max size in bytes, for instance 'limit:1k')"),
		withConfigBuilder(limitConfigBuilder))
}

func limit(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	bytes := config.(int64)

	log.Debugf("limiting to %d bytes\n", bytes)
	lr := io.LimitReader(r, bytes)

	return io.Copy(w, lr) // streamable
}

func limitConfigBuilder(args []string) (any, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumberOfArgs(1, len(args))
	}

	n, err := units.FromHumanSize(args[0])
	if err != nil {
		return nil, err
	}
	return n, nil
}
