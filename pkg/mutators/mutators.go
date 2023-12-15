package mutators

import (
	"fmt"
	"io"
	"os"
	"slices"
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

const (
	argSeparator = ":"
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
	NewMutator(logger *log.Logger, args []string) (Mutator, error)
	Name() string
	Description() string
	Category() string
}

type mutatorCollection struct {
	factories map[string]Factory
	aliases   map[string]string
	logger    *log.Logger
	Name      string
	mutators  []Mutator
	mu        sync.Mutex
}

func newCollection(name string, logger *log.Logger) *mutatorCollection {
	glog = logger
	defer glog.Printf("mutators: collection %s ready.\n", name)

	return &mutatorCollection{
		Name:      name,
		factories: make(map[string]Factory),
		aliases:   make(map[string]string),
		logger:    logger,
	}
}

func RegisterFactory(name string, factory Factory) error {
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

func RegisterAlias(name string, alias string) error {
	globalCollection.mu.Lock()
	defer globalCollection.mu.Unlock()
	// forbid overwriting an existing alias
	if _, ok := globalCollection.aliases[alias]; ok {
		return fmt.Errorf("mutators: aliasing %s is not permitted. It is already an alias", name)
	}
	globalCollection.aliases[alias] = name
	// glog.Printf("mutators: %s aliased as %s\n", name, alias)
	return nil
}

func New(fullName string) (Mutator, error) {
	globalCollection.mu.Lock()
	defer globalCollection.mu.Unlock()

	name, argsString, argsFound := strings.Cut(fullName, argSeparator)

	var args []string
	if argsString != "" {
		args = strings.Split(argsString, argSeparator)
	}

	if factoryName, ok := globalCollection.aliases[name]; ok {
		glog.Printf("mutators: %s is an alias to %s\n", name, factoryName)
		name = factoryName
	}
	factory, ok := globalCollection.factories[name]
	if !ok {
		TryFuzzySearch(name)
		return nil, fmt.Errorf("mutators: %s not found", name)
	}
	glog.Printf("mutators: instancing %s\n", name)
	if argsFound {
		glog.Printf("mutators: with args %v\n", args)
	}

	m, err := factory.NewMutator(globalCollection.logger, args)
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

func ListAvailableAliases() ([]string, map[string][]string) {
	detailed := make(map[string][]string)
	aliases := make([]string, 0, len(globalCollection.aliases))
	for alias, factory := range globalCollection.aliases {
		aliases = append(aliases, alias)
		detailed[factory] = append(detailed[factory], alias)
	}
	return aliases, detailed
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
	s.WriteString("\n  ('X:Y' means X is an argument with default value Y)\n")
	_, d := ListAvailableAliases()
	if len(d) > 0 {
		s.WriteString("\n  mutator aliases:\n")
	}
	// for factory, aliases := range d {
	// 	fmt.Fprintf(&s, "    %s: %s\n", strings.Join(aliases, ", "), factory)
	// }
	keys = make([]string, 0, len(d))
	for k := range d {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, factory := range keys {
		fmt.Fprintf(&s, "    %s: %s\n", strings.Join(d[factory], ", "), factory)
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
	list, _ := ListAvailableAliases()
	list = append(list, ListAvailableMutators("ALL")...)
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
