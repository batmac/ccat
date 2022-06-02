package mutator

import (
	"ccat/log"
	"fmt"
	"io"
)

var name = "dummy"
var description = "simple fifo"

type dummyMutator struct {
	genericMutator
}

func init() {
	register(name, newDummy)
}

func newDummy(logger *log.Logger) (Mutator, error) {
	logger.Println("dummy: new")
	return &dummyMutator{
		genericMutator: newGeneric(logger),
	}, nil
}

func (m *dummyMutator) Start(w io.WriteCloser, r io.ReadCloser) error {
	m.mu.Lock()
	if m.started {
		m.mu.Unlock()
		return fmt.Errorf("dummy: mutator has already started.")
	}
	m.started = true
	m.mu.Unlock()
	m.logger.Printf("dummy: start\n")

	go func() {
		m.logger.Printf("dummy: copying\n")
		written, err := io.Copy(w, r)
		m.logger.Printf("dummy: done\n")
		if err != nil {
			m.logger.Println(err)
		}
		m.logger.Printf("dummy: written %d bytes\n", written)
		m.logger.Printf("dummy: closing %v\n", w)
		w.Close()
		if err != nil {
			m.logger.Println(err)
		}
		close(m.done)
	}()

	return nil
}
func (m *dummyMutator) Wait() error {
	m.logger.Printf("dummy: wait called\n")
	m.mu.Lock()
	if !m.started {
		m.mu.Unlock()
		return fmt.Errorf("dummy: mutator is not started")
	}
	if m.waited {
		m.mu.Unlock()
		return fmt.Errorf("dummy: mutator is already waited")
	}
	m.waited = true
	m.mu.Unlock()
	<-m.done
	return nil
}

func (m *dummyMutator) Name() string {
	return name
}
func (m *dummyMutator) Description() string {
	return description
}
