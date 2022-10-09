package mutators

import (
	"sync"

	"github.com/batmac/ccat/pkg/log"
)

type GenericMutator struct {
	Logger  *log.Logger
	Done    chan struct{}
	Mu      sync.Mutex
	Started bool
	Waited  bool
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
