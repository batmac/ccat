package globalctx

import (
	"ccat/log"
	"context"
	"sync"
)

var globalCtx ctx

type ctx struct {
	ctx context.Context
	mu  sync.Mutex
}

func init() {
	globalCtx = ctx{ctx: context.Background()}
}

func Set(k string, v interface{}) {
	globalCtx.mu.Lock()
	defer globalCtx.mu.Unlock()
	log.Debugf("globalctx: setting %v=%v\n", k, v)
	globalCtx.ctx = context.WithValue(globalCtx.ctx, k, v)
}

func Get(k string) interface{} {
	globalCtx.mu.Lock()
	defer globalCtx.mu.Unlock()
	return globalCtx.ctx.Value(k)
}
