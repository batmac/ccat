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
	globalCollection = NewCollection("globalMutatorsCollection", log.Default())
)

type Mutator interface {
	Start(w io.WriteCloser, r io.ReadCloser) error
	Wait() error
	Name() string
	Description() string
}
type Factory interface {
	New(logger *log.Logger) (Mutator, error)
	Name() string
	Description() string
}

type MutatorCollection struct {
	mu       sync.Mutex
	Name     string
	mutators []Mutator
	//Mutators  map[string]Mutator
	factories map[string]Factory
	logger    *log.Logger
}

func NewCollection(name string, logger *log.Logger) *MutatorCollection {

	glog = logger
	defer glog.Printf("mutators: collection %s ready.\n", name)

	return &MutatorCollection{
		Name: name,
		//Mutators:  make(map[string]Mutator),
		factories: make(map[string]Factory),
		logger:    logger,
	}
}

func register(name string, factory Factory) error {
	globalCollection.mu.Lock()
	if _, ok := globalCollection.factories[name]; ok {
		return fmt.Errorf("mutators: %s is already registered", name)
	}
	globalCollection.factories[name] = factory
	globalCollection.mu.Unlock()
	glog.Printf("mutators: %s registered\n", name)
	return nil
}

func RunTest(name string, w io.WriteCloser, r io.ReadCloser) error {
	globalCollection.mu.Lock()
	defer globalCollection.mu.Unlock()

	if factory, ok := globalCollection.factories[name]; !ok {
		return fmt.Errorf("mutator %s not found", name)
	} else {
		glog.Printf(" instancing mutator %s\n", name)

		m, err := factory.New(globalCollection.logger)
		if err != nil {
			return err
		}
		globalCollection.mutators = append(globalCollection.mutators, m)
		glog.Printf(" instanced mutator %s\n", name)

		glog.Printf(" starting mutator %s\n", name)

		err = m.Start(w, r)
		if err != nil {
			return err
		}
		glog.Printf(" waiting mutator %s\n", name)

		err = m.Wait()
		if err != nil {
			return err
		}
	}
	glog.Printf(" returning from runtest")

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

	m, err := factory.New(globalCollection.logger)
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
