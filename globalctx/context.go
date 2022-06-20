package globalctx

import (
	"context"
	"sync"

	"github.com/batmac/ccat/log"
)

// globalctx is used to set/get context to the processing pipeline
// (typically set by mutators)
// keys could be:
// fileList: the file arguments
// path: path (url) of the current processed file
// hintLexer: hint about the lexer the highlighter should probably use
// expectingBinary: the output will have non-displayable char(so don't try to pretty-print/highlight)
var globalCtx ctx

type ctx struct {
	ctx context.Context
	mu  sync.Mutex
}

type key string

func init() {
	Reset()
}

func Set(k string, v interface{}) {
	globalCtx.mu.Lock()
	defer globalCtx.mu.Unlock()
	log.Debugf("globalctx: setting %v=%v\n", k, v)
	globalCtx.ctx = context.WithValue(globalCtx.ctx, key(k), v)
}

func Get(k string) interface{} {
	globalCtx.mu.Lock()
	defer globalCtx.mu.Unlock()
	return globalCtx.ctx.Value(key(k))
}

func Reset() {
	globalCtx.mu.Lock()
	defer globalCtx.mu.Unlock()
	globalCtx.ctx = context.Background()
}
