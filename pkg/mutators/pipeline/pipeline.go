package pipeline

import (
	"errors"
	"io"
	"strings"
	"sync"

	"github.com/batmac/ccat/pkg/config"
	"github.com/batmac/ccat/pkg/globalctx"
	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/mutators"
)

var globalPipeline pipeline

type pipeline struct {
	stages []mutators.Mutator
	mu     sync.Mutex
}

func NewPipeline(description string, out io.WriteCloser, in io.ReadCloser) error {
	globalPipeline.mu.Lock()
	if len(globalPipeline.stages) > 0 {
		log.Fatal("pipeline is not empty\n")
	}
	if len(description) == 0 {
		globalPipeline.mu.Unlock()
		return errors.New("empty pipeline requested")
	}
	description = handleAliases(description)
	list := strings.Split(description, ",")
	for _, m := range list {
		log.Debugf("creating %v\n", m)
		mutator, err := mutators.New(m)
		if err != nil {
			log.Fatalf("mutator '%s': %s\n", m, err.Error())
		}
		globalPipeline.stages = append(globalPipeline.stages, mutator)
	}
	globalPipeline.mu.Unlock()

	ready := make(chan struct{})
	go func() {
		globalPipeline.mu.Lock()
		from, to := in, out
		for _, mutator := range globalPipeline.stages {
			r, w := io.Pipe()
			log.Debugf("starting %v\n", mutator.Name())
			if mutator.Start(w, from) != nil {
				log.Fatal("failed to start the mutator\n")
			}
			from = r
		}
		close(ready) // all mutators are started, we are ready
		globalPipeline.mu.Unlock()
		n, err := io.Copy(to, from)
		if err != nil {
			log.Fatal(err)
		}
		log.Debugf("copied %v bytes.", n)
		log.Debugf("closing pipeline.\n")
		err = from.Close()
		if err != nil {
			log.Debugln(err)
		}
		err = to.Close()
		if err != nil {
			log.Debugln(err)
		}
		log.Debugf("closed pipeline.\n")
	}()
	<-ready
	return nil
}

func Wait() {
	globalPipeline.mu.Lock()
	defer globalPipeline.mu.Unlock()
	for _, m := range globalPipeline.stages {
		log.Debugf("waiting %v\n", m)

		err := m.Wait()
		if err != nil {
			//nolint:gocritic // exitAfterDefer
			log.Fatal(err)
		}
		log.Debugf("waited %v\n", m)
	}
	globalPipeline.stages = nil
}

func handleAliases(description string) string {
	c, ok := globalctx.Get("conf").(*config.Config)

	if !ok || len(c.Aliases) == 0 {
		// no aliases defined
		return description
	}

	for alias, definition := range c.Aliases {
		alias = "@" + alias
		description = strings.ReplaceAll(description, alias, definition)
	}

	log.Debugf("pipeline definition => %v\n", description)

	if definition, alias, found := strings.Cut(description, "="); found {
		for alias[0] == '@' {
			alias = alias[1:]
		}
		c.Aliases[alias] = definition
		globalctx.Set("conf", c)
		log.Debugf("added alias %v=%v\n", alias, definition)
		return definition
	}
	return description
}
