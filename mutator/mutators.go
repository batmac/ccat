package mutator

import (
	"ccat/log"
	"fmt"
	"io"
	"sync"
)

var (
	// register() is called from init() so this has to be global
	glog             *log.Logger
	globalCollection = NewCollection("globalCollection", log.Default())
)

type Mutator interface {
	Start(w io.WriteCloser, r io.ReadCloser) error
	Wait() error
	Name() string
	Description() string
}
type Factory func(logger *log.Logger) (Mutator, error)

type MutatorCollection struct {
	sync.Mutex
	Name     string
	Mutators []Mutator
	//Mutators  map[string]Mutator
	factories map[string]Factory
	logger    *log.Logger
}

func NewCollection(name string, logger *log.Logger) *MutatorCollection {

	glog = logger
	glog.Printf("mutators collection %s ready.\n", name)

	return &MutatorCollection{
		Name: name,
		//Mutators:  make(map[string]Mutator),
		factories: make(map[string]Factory),
		logger:    logger,
	}
}

func register(name string, factory Factory) error {
	globalCollection.Lock()
	if _, ok := globalCollection.factories[name]; ok {
		return fmt.Errorf("%s is already registered", name)
	}
	globalCollection.factories[name] = factory
	globalCollection.Unlock()
	glog.Printf(" mutator %s registered\n", name)
	return nil
}

func RunTest(name string, w io.WriteCloser, r io.ReadCloser) error {
	globalCollection.Lock()
	defer globalCollection.Unlock()

	if f, ok := globalCollection.factories[name]; !ok {
		return fmt.Errorf("mutator %s not found", name)
	} else {
		glog.Printf(" instancing mutator %s\n", name)

		m, err := f(globalCollection.logger)
		if err != nil {
			return err
		}
		globalCollection.Mutators = append(globalCollection.Mutators, m)
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
