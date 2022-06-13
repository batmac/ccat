package mutators

import (
	"fmt"
	"io"

	"github.com/batmac/ccat/globalctx"
	"github.com/batmac/ccat/log"
)

// launch a mutator in its dedicated goroutine

type simpleFn func(w io.WriteCloser, r io.ReadCloser) (int64, error)

type simpleMutator struct {
	genericMutator
	factory *simpleFactory
}

type simpleFactory struct {
	name, description string
	fn                simpleFn
	hintLexer         string
	expectingBinary   bool
}

type simpleOption func(*simpleFactory)

func withHintLexer(s string) simpleOption {
	return func(f *simpleFactory) {
		f.hintLexer = s
	}
}
func withDescription(s string) simpleOption {
	return func(f *simpleFactory) {
		f.description = s
	}
}
func withExpectingBinary(b bool) simpleOption {
	return func(f *simpleFactory) {
		f.expectingBinary = true
	}
}

func simpleRegister(name string, f simpleFn, opts ...simpleOption) {
	factory := new(simpleFactory)
	factory.name = name
	factory.fn = f
	for _, o := range opts {
		o(factory)
	}
	register(name, factory)
}

func (f *simpleFactory) newMutator(logger *log.Logger) (Mutator, error) {
	logger.Printf("%s: new", f.Name())
	globalctx.Set("hintLexer", f.hintLexer)
	globalctx.Set("expectingBinary", f.expectingBinary)

	return &simpleMutator{
		genericMutator: newGeneric(logger),
		factory:        f,
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
		written, err := m.factory.fn(w, r)
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
	return m.factory.Name()
}
func (m *simpleMutator) Description() string {
	return m.factory.Description()
}
func (f *simpleFactory) Name() string {
	return f.name
}
func (f *simpleFactory) Description() string {
	return f.description
}
