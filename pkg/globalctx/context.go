package globalctx

import (
	"context"
	"sync"

	"github.com/batmac/ccat/pkg/log"
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
	err bool
}

type key string

func init() {
	globalCtx.mu.Lock()
	defer globalCtx.mu.Unlock()
	globalCtx.ctx = context.Background()
}

func Set(k string, v any) {
	globalCtx.mu.Lock()
	defer globalCtx.mu.Unlock()
	log.Debugf("globalctx: setting %v=%#v\n", k, v)
	globalCtx.ctx = context.WithValue(globalCtx.ctx, key(k), v)
}

func Get(k string) any {
	globalCtx.mu.Lock()
	defer globalCtx.mu.Unlock()
	return globalCtx.ctx.Value(key(k))
}

func GetBool(k string) bool {
	globalCtx.mu.Lock()
	defer globalCtx.mu.Unlock()
	bi := globalCtx.ctx.Value(key(k))
	if b, ok := bi.(bool); ok {
		return b
	}
	return false
}

func Reset() {
	// doesn't touch globalCtx.err, copy the conf key
	globalCtx.mu.Lock()
	defer globalCtx.mu.Unlock()
	conf := globalCtx.ctx.Value(key("conf"))
	globalCtx.ctx = context.Background()
	globalCtx.ctx = context.WithValue(globalCtx.ctx, key("conf"), conf)
}

func SetErrored() {
	globalCtx.mu.Lock()
	defer globalCtx.mu.Unlock()
	globalCtx.err = true
}

func IsErrored() bool {
	globalCtx.mu.Lock()
	defer globalCtx.mu.Unlock()
	return globalCtx.err
}
