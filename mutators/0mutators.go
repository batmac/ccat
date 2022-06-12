package mutators

import (
	"ccat/log"
	"fmt"
	"io"
	"sort"
	"sync"
)

var (
	// register() is called from init() so this has to be global
	glog             *log.Logger // shortcut for globalCollection.logger
	globalCollection = newCollection("globalMutatorsCollection", log.Default())
)

//Mutator and factory should be totally separate or reentrant as in the future they may be used simultanetely

type Mutator interface {
	Start(w io.WriteCloser, r io.ReadCloser) error
	Wait() error
	Name() string
	Description() string
}
type factory interface {
	new(logger *log.Logger) (Mutator, error)
	Name() string
	Description() string
}

type mutatorCollection struct {
	mu       sync.Mutex
	Name     string
	mutators []Mutator
	//Mutators  map[string]Mutator
	factories map[string]factory
	logger    *log.Logger
}

func newCollection(name string, logger *log.Logger) *mutatorCollection {

	glog = logger
	defer glog.Printf("mutators: collection %s ready.\n", name)

	return &mutatorCollection{
		Name: name,
		//Mutators:  make(map[string]Mutator),
		factories: make(map[string]factory),
		logger:    logger,
	}
}

func register(name string, factory factory) error {
	globalCollection.mu.Lock()
	if _, ok := globalCollection.factories[name]; ok {
		return fmt.Errorf("mutators: %s is already registered", name)
	}
	globalCollection.factories[name] = factory
	globalCollection.mu.Unlock()
	glog.Printf("mutators: %s registered\n", name)
	return nil
}

func New(name string) (Mutator, error) {
	globalCollection.mu.Lock()
	defer globalCollection.mu.Unlock()

	factory, ok := globalCollection.factories[name]
	if !ok {
		return nil, fmt.Errorf("mutators: %s not found", name)
	}
	glog.Printf("mutators: instancing %s\n", name)

	m, err := factory.new(globalCollection.logger)
	if err != nil {
		return nil, err
	}
	globalCollection.mutators = append(globalCollection.mutators, m)
	glog.Printf("mutators: returning a new %s\n", name)
	return m, nil
}

func ListMutators() []string {
	var l []string
	for _, v := range globalCollection.factories {
		l = append(l, v.Name()+": "+v.Description())
	}
	sort.Strings(l)
	return l
}
