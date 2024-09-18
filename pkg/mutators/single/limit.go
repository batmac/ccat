package mutators

import (
	"io"

	"github.com/batmac/ccat/pkg/log"
)

func init() {
	singleRegister("limit", limit, withDescription("a simple limiting fifo ( with X max size in bytes, for instance 'limit:1k')"),
		withConfigBuilder(stdConfigHumanSizeAsInt64WithDefault(500)),
		withAliases("l", "head"))
}

func limit(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	bytes := config.(int64)

	log.Debugf("limiting to %d bytes\n", bytes)
	lr := io.LimitReader(r, bytes)

	return io.Copy(w, lr) // streamable
}
