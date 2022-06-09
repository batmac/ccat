package globalctx

import (
	"ccat/log"
	"context"
)

var globalCtx context.Context

func init() {
	globalCtx = context.Background()
}

func Set(k string, v interface{}) {
	log.Debugf("globalctx: setting %v=%v\n", k, v)
	globalCtx = context.WithValue(globalCtx, k, v)
}

func Get(k string) interface{} {
	return globalCtx.Value(k)
}
