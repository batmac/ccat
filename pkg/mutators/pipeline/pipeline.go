package pipeline

import (
	"errors"
	"io"
	"strings"
	"sync"

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
	list := strings.Split(description, ",")
	for _, m := range list {
		log.Debugf("creating %v\n", m)
		mutator, err := mutators.New(m)
		if err != nil {
			log.Fatal(err)
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
