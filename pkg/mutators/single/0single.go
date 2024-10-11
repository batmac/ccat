package mutators

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/batmac/ccat/pkg/globalctx"
	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/mutators"
)

// launch a mutator in its dedicated goroutine

type (
	configBuilder func(args []string) (any, error)
	singleOption  func(*singleFactory)
	singleFn      func(w io.WriteCloser, r io.ReadCloser, config any) (int64, error)
)

type singleMutator struct {
	config  any
	factory *singleFactory
	Logger  *log.Logger
	Done    chan struct{}
	Mu      sync.Mutex
	Started bool
	Waited  bool
}

type singleFactory struct {
	fn                singleFn
	configBuilder     configBuilder
	name, description string
	category          string
	hintLexer         string
	hintSlowOutput    bool
	expectingBinary   bool
}

func ErrWrongNumberOfArgs(amin, amax, got int) error {
	if amin == amax {
		return fmt.Errorf("wrong number of arguments, got %d, expected %d", got, amin)
	}
	return fmt.Errorf("wrong number of arguments, expected between %d and %d, got %d", amin, amax, got)
}

func withHintLexer(s string) singleOption {
	return func(f *singleFactory) {
		f.hintLexer = s
	}
}

func withDescription(s string) singleOption {
	return func(f *singleFactory) {
		f.description = s
	}
}

func withCategory(s string) singleOption {
	return func(f *singleFactory) {
		if s == "compress" {
			f.expectingBinary = true
		}
		f.category = s
	}
}

func withExpectingBinary() singleOption {
	return func(f *singleFactory) {
		f.expectingBinary = true
	}
}

func withHintSlow() singleOption {
	return func(f *singleFactory) {
		f.hintSlowOutput = true
		// if slow, we set expectingBinary to true as well not to highlight
		f.expectingBinary = true
	}
}

func withConfigBuilder(fn configBuilder) singleOption {
	return func(f *singleFactory) {
		f.configBuilder = fn
	}
}

func withAliases(aliases ...string) singleOption {
	return func(f *singleFactory) {
		for _, alias := range aliases {
			_ = mutators.RegisterAlias(f.name, alias)
		}
	}
}

func defaultConfigBuilder(args []string) (any, error) {
	// no config, no args permitted
	if len(args) != 0 {
		return nil, ErrWrongNumberOfArgs(0, 0, len(args))
	}
	return nil, nil //nolint:nilnil
}

func singleRegister(name string, f singleFn, opts ...singleOption) {
	factory := new(singleFactory)
	factory.name = name
	factory.fn = f
	factory.configBuilder = defaultConfigBuilder
	for _, o := range opts {
		o(factory)
	}
	if err := mutators.RegisterFactory(name, factory); err != nil {
		log.Debugf("registering %s failed!\n", name)
	}
}

func (f *singleFactory) NewMutator(logger *log.Logger, args []string) (mutators.Mutator, error) {
	logger.Printf("%s: new", f.Name())
	globalctx.Set("hintLexer", f.hintLexer)
	globalctx.Set("hintSlowOutput", f.hintSlowOutput)
	globalctx.Set("expectingBinary", f.expectingBinary)

	var config any
	var err error
	if f.configBuilder != nil {
		config, err = f.configBuilder(args)
		if err != nil {
			return nil, err
		}
	}

	return &singleMutator{
		Logger:  logger,
		Done:    make(chan struct{}),
		factory: f,
		config:  config,
	}, nil
}

func (m *singleMutator) Start(w io.WriteCloser, r io.ReadCloser) error {
	m.Mu.Lock()
	if m.Started {
		m.Mu.Unlock()
		return fmt.Errorf("%s: mutator has already started", m.Name())
	}
	m.Started = true
	m.Mu.Unlock()
	m.Logger.Printf("%s: start %v\n", m.Name(), w)

	go func() {
		m.Logger.Printf("%s: dumping from %v to %v\n", m.Name(), r, w)
		written, err := m.factory.fn(w, r, m.config)
		if err != nil && !errors.Is(err, io.ErrClosedPipe) {
			log.Fatal(err)
		}
		m.Logger.Printf("%s: done\n", m.Name())
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

func (m *singleMutator) Wait() error {
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

func (m *singleMutator) Name() string {
	return m.factory.Name()
}

func (m *singleMutator) Description() string {
	return m.factory.Description()
}

func (m *singleMutator) Category() string {
	return m.factory.Category()
}

func (f *singleFactory) Name() string {
	return f.name
}

func (f *singleFactory) Description() string {
	return f.description
}

func (f *singleFactory) Category() string {
	return f.category
}
