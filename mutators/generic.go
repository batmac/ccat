package mutators

import (
	"sync"

	"github.com/batmac/ccat/log"
)

type GenericMutator struct {
	Mu     sync.Mutex
	Logger *log.Logger

	Started bool
	Waited  bool
	Done    chan struct{}
}

func NewGeneric(logger *log.Logger) GenericMutator {
	return GenericMutator{
		Mu:     sync.Mutex{},
		Logger: logger,

		Started: false,
		Waited:  false,
		Done:    make(chan struct{}),
	}
}
