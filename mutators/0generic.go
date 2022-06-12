package mutators

import (
	"ccat/log"
	"sync"
)

type genericMutator struct {
	mu     sync.Mutex
	logger *log.Logger

	started bool
	waited  bool
	done    chan struct{}
}

func newGeneric(logger *log.Logger) genericMutator {
	return genericMutator{
		mu:     sync.Mutex{},
		logger: logger,

		started: false,
		waited:  false,
		done:    make(chan struct{}),
	}
}
