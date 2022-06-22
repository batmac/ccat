package pipeline

import (
	"io"
	"strings"

	"github.com/batmac/ccat/log"
	"github.com/batmac/ccat/mutators"
)

var globalPipeline []mutators.Mutator

func NewPipeline(description string, out io.WriteCloser, in io.ReadCloser) error {
	list := strings.Split(description, ",")
	for _, m := range list {
		log.Debugf("creating %v\n", m)
		mutator, err := mutators.New(m)
		if err != nil {
			log.Fatal(err)
		}
		globalPipeline = append(globalPipeline, mutator)
	}
	go func() {
		from, to := in, out
		for _, mutator := range globalPipeline {
			r, w := io.Pipe()
			log.Debugf("starting %v\n", mutator.Name())
			if mutator.Start(w, from) != nil {
				log.Fatal("failed to start the mutator\n")
			}
			from = r
		}
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

	return nil
}

func Wait() {
	for _, m := range globalPipeline {
		log.Debugf("waiting %v\n", m)

		err := m.Wait()
		if err != nil {
			log.Fatal(err)
		}
		log.Debugf("waited %v\n", m)
	}
}

func Reset() {
	globalPipeline = nil
}
