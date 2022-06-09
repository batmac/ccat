package mutators

import (
	"ccat/globalctx"
	"ccat/log"
	"fmt"
	"io"
)

type simpleFn func(w io.WriteCloser, r io.ReadCloser) (int64, error)

type simpleMutator struct {
	genericMutator
	name, description string
	fn                simpleFn
	hintLexer         string
}

type simpleFactory struct {
	name, description string
	fn                simpleFn
	hintLexer         string
}

func simpleRegister(name, description, hintLexer string, f simpleFn) {
	factory := new(simpleFactory)
	factory.name = name
	factory.description = description
	factory.fn = f
	factory.hintLexer = hintLexer

	register(name, factory)
}

func (f *simpleFactory) New(logger *log.Logger) (Mutator, error) {
	logger.Printf("%s: new", f.Name())
	if len(f.hintLexer) != 0 {
		globalctx.Set("hintLexer", f.hintLexer)
	}
	return &simpleMutator{
		genericMutator: newGeneric(logger),
		name:           f.name,
		description:    f.description,
		fn:             f.fn,
		hintLexer:      f.hintLexer,
	}, nil
}

func (m *simpleMutator) Start(w io.WriteCloser, r io.ReadCloser) error {
	m.mu.Lock()
	if m.started {
		m.mu.Unlock()
		return fmt.Errorf("%s: mutator has already started.", m.Name())
	}
	m.started = true
	m.mu.Unlock()
	m.logger.Printf("%s: start %v\n", m.Name(), w)

	go func() {
		m.logger.Printf("%s: dumping from %v to %v\n", m.Name(), r, w)
		written, err := m.fn(w, r)
		m.logger.Printf("%s: done\n", m.Name())
		if err != nil {
			log.Fatal(err)
		}
		m.logger.Printf("%s: written %d bytes\n", m.Name(), written)
		m.logger.Printf("%s: closing %v\n", m.Name(), w)
		w.Close()
		if err != nil {
			m.logger.Println(err)
		}
		close(m.done)
	}()

	return nil
}
func (m *simpleMutator) Wait() error {
	m.logger.Printf("%s: wait called\n", m.Name())
	m.mu.Lock()
	if !m.started {
		m.mu.Unlock()
		return fmt.Errorf("%s: mutator is not started", m.Name())
	}
	if m.waited {
		m.mu.Unlock()
		return fmt.Errorf("%s: mutator is already waited", m.Name())
	}
	m.waited = true
	m.mu.Unlock()
	<-m.done
	return nil
}

func (m *simpleMutator) Name() string {
	return m.name
}
func (m *simpleMutator) Description() string {
	return m.description
}
func (f *simpleFactory) Name() string {
	return f.name
}
func (f *simpleFactory) Description() string {
	return f.description
}
