package mutators

import (
	"ccat/log"
	"fmt"
	"io"
)

var dummyName = "dummy"
var dummyDescription = "simple fifo"

type dummyMutator struct {
	genericMutator
}

type dummyFactory struct {
}

func init() {
	f := new(dummyFactory)
	register(dummyName, f)
}

func (f *dummyFactory) New(logger *log.Logger) (Mutator, error) {
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
	m.logger.Printf("dummy: start %v\n", w)

	go func() {
		m.logger.Printf("dummy: copying from %v to %v\n", r, w)
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
	return dummyName
}
func (m *dummyMutator) Description() string {
	return dummyDescription
}
func (f *dummyFactory) Name() string {
	return dummyName
}
func (f *dummyFactory) Description() string {
	return dummyDescription
}
