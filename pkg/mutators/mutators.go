package mutators

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/stringutils"
)

var (
	// register() is called from init() so this has to be global
	glog             *log.Logger // shortcut for globalCollection.logger
	globalCollection = newCollection("globalMutatorsCollection", log.Default())
)

// Mutator and factory should be totally separate or reentrant as they may be used simultaneously

type Mutator interface {
	Start(w io.WriteCloser, r io.ReadCloser) error
	Wait() error
	Name() string
	Description() string
	Category() string
}
type Factory interface {
	NewMutator(logger *log.Logger) (Mutator, error)
	Name() string
	Description() string
	Category() string
}

type mutatorCollection struct {
	mu       sync.Mutex
	Name     string
	mutators []Mutator
	// Mutators  map[string]Mutator
	factories map[string]Factory
	logger    *log.Logger
}

func newCollection(name string, logger *log.Logger) *mutatorCollection {
	glog = logger
	defer glog.Printf("mutators: collection %s ready.\n", name)

	return &mutatorCollection{
		Name: name,
		// Mutators:  make(map[string]Mutator),
		factories: make(map[string]Factory),
		logger:    logger,
	}
}

func Register(name string, factory Factory) error {
	globalCollection.mu.Lock()
	if _, ok := globalCollection.factories[name]; ok {
		globalCollection.mu.Unlock()
		return fmt.Errorf("mutators: %s is already registered", name)
	}
	globalCollection.factories[name] = factory
	globalCollection.mu.Unlock()
	// glog.Printf("mutators: %s registered\n", name)
	return nil
}

func New(name string) (Mutator, error) {
	globalCollection.mu.Lock()
	defer globalCollection.mu.Unlock()

	factory, ok := globalCollection.factories[name]
	if !ok {
		TryFuzzySearch(name)
		return nil, fmt.Errorf("mutators: %s not found", name)
	}
	glog.Printf("mutators: instancing %s\n", name)

	m, err := factory.NewMutator(globalCollection.logger)
	if err != nil {
		return nil, err
	}
	globalCollection.mutators = append(globalCollection.mutators, m)
	glog.Printf("mutators: returning a new %s\n", name)
	return m, nil
}

func ListAvailableMutators(category string) []string {
	l := make([]string, 0, len(globalCollection.factories))
	for _, v := range globalCollection.factories {
		if category == "ALL" || category == v.Category() {
			l = append(l, v.Name())
		}
	}
	sort.Strings(l)
	return l
}

func listAvailableMutatorsByCategoryWithDescriptions() map[string][]string {
	listByCategory := make(map[string][]string)
	for _, v := range globalCollection.factories {
		listByCategory[v.Category()] = append(listByCategory[v.Category()], v.Name()+": "+v.Description())
	}
	return listByCategory
}

func AvailableMutatorsHelp() string {
	var s strings.Builder
	l := listAvailableMutatorsByCategoryWithDescriptions()
	keys := make([]string, 0, len(l))
	for k := range l {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, category := range keys {
		if len(category) > 0 {
			s.WriteString("    " + category + ":\n")
		}
		sort.Strings(l[category])
		for _, mutator := range l[category] {
			s.WriteString("        " + mutator + "\n")
		}
	}
	return s.String()
}

func Run(mutatorName, input string) string {
	in := io.NopCloser(strings.NewReader(input))
	out := new(stringutils.NopStringWriteCloser)
	m, err := New(mutatorName)
	if err != nil {
		log.Fatal(err)
	}
	if m.Start(out, in) != nil {
		log.Fatal("failed to start the mutator\n")
	}
	err = m.Wait()
	if err != nil {
		log.Fatal(err)
	}
	return out.String()
}

func TryFuzzySearch(name string) {
	list := ListAvailableMutators("ALL")
	f, err := stringutils.FuzzySearch(name, list, 0.5)
	if err != nil {
		log.Debugln(err)
		return
	}
	if len(f) == 0 {
		log.Debugf("fuzzysearch found nothing for '%s'\n", name)
		return
	}
	fmt.Fprintf(os.Stderr, "'%s' does not exist, did you mean %s ?\n", name, f)
	os.Exit(1)
}
