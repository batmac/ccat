package mutators

import (
	"io"

	"github.com/batmac/ccat/pkg/log"

	"github.com/docker/go-units"
)

func init() {
	simpleRegister("limit", limit, withDescription("a simple limiting fifo (max size in bytes, for instance 'limit:1k')"),
		withArgsValidator(limitValidator))
}

func limit(w io.WriteCloser, r io.ReadCloser, args ...string) (int64, error) {
	bytes, err := units.FromHumanSize(args[0])
	if err != nil {
		return 0, err
	}
	log.Debugf("limiting to %d bytes\n", bytes)
	lr := io.LimitReader(r, bytes)

	return io.Copy(w, lr) // streamable
}

func limitValidator(args []string) error {
	if len(args) != 1 {
		return ErrWrongNumberOfArgs(1, len(args))
	}

	if _, err := units.FromHumanSize(args[0]); err != nil {
		return err
	}
	/* 	if _, err := strconv.ParseInt(args[0], 10, 64); err != nil {
		return err
	} */
	return nil
}
