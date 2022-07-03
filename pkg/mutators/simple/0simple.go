package mutators

import (
	"fmt"
	"io"

	"github.com/batmac/ccat/pkg/globalctx"
	"github.com/batmac/ccat/pkg/log"
	. "github.com/batmac/ccat/pkg/mutators"
)

// launch a mutator in its dedicated goroutine

type simpleFn func(w io.WriteCloser, r io.ReadCloser) (int64, error)

type simpleMutator struct {
	GenericMutator
	factory *simpleFactory
}

type simpleFactory struct {
	name, description string
	category          string
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

func withCategory(s string) simpleOption {
	return func(f *simpleFactory) {
		if s == "compress" {
			f.expectingBinary = true
		}
		f.category = s
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
	if err := Register(name, factory); err != nil {
		log.Debugf("registering %s failed!\n", name)
	}
}

func (f *simpleFactory) NewMutator(logger *log.Logger) (Mutator, error) {
	logger.Printf("%s: new", f.Name())
	globalctx.Set("hintLexer", f.hintLexer)
	globalctx.Set("expectingBinary", f.expectingBinary)

	return &simpleMutator{
		GenericMutator: NewGeneric(logger),
		factory:        f,
	}, nil
}

func (m *simpleMutator) Start(w io.WriteCloser, r io.ReadCloser) error {
	m.Mu.Lock()
	if m.Started {
		m.Mu.Unlock()
		return fmt.Errorf("%s: mutator has already started.", m.Name())
	}
	m.Started = true
	m.Mu.Unlock()
	m.Logger.Printf("%s: start %v\n", m.Name(), w)

	go func() {
		m.Logger.Printf("%s: dumping from %v to %v\n", m.Name(), r, w)
		written, err := m.factory.fn(w, r)
		m.Logger.Printf("%s: done\n", m.Name())
		if err != nil {
			log.Fatal(err)
		}
		m.Logger.Printf("%s: written %d bytes\n", m.Name(), written)
		m.Logger.Printf("%s: closing %v\n", m.Name(), w)
		err = r.Close()
		if err != nil {
			m.Logger.Println(err)
		}
		err = w.Close()
		if err != nil {
			m.Logger.Println(err)
		}
		m.Logger.Printf("%s: closed %v\n", m.Name(), w)
		close(m.Done)
	}()

	return nil
}

func (m *simpleMutator) Wait() error {
	m.Logger.Printf("%s: wait called\n", m.Name())
	m.Mu.Lock()
	if !m.Started {
		m.Mu.Unlock()
		return fmt.Errorf("%s: mutator is not started", m.Name())
	}
	if m.Waited {
		m.Mu.Unlock()
		return fmt.Errorf("%s: mutator is already waited", m.Name())
	}
	m.Waited = true
	m.Mu.Unlock()
	<-m.Done
	return nil
}

func (m *simpleMutator) Name() string {
	return m.factory.Name()
}

func (m *simpleMutator) Description() string {
	return m.factory.Description()
}

func (m *simpleMutator) Category() string {
	return m.factory.Category()
}

func (f *simpleFactory) Name() string {
	return f.name
}

func (f *simpleFactory) Description() string {
	return f.description
}

func (f *simpleFactory) Category() string {
	return f.category
}
