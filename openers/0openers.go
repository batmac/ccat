package openers

import (
	"ccat/log"
	"errors"
	"io"
	"io/fs"
	"os"
	"strings"
	"sync"
)

var (
	// register() is called from init() so this has to be global
	globalCollection = newCollection("global")
)

type Opener interface {
	Open(s string, lock bool) (io.ReadCloser, error)
	Evaluate(s string) float32 //score the ability to open
	Name() string
	Description() string
}

type OpenerCollection struct {
	mu      sync.Mutex
	Name    string
	openers []Opener
}

func newCollection(name string) *OpenerCollection {
	//log.Printf("openers collection %s ready.\n", name)
	return &OpenerCollection{
		Name: name,
	}
}

func register(opener Opener) error {
	globalCollection.mu.Lock()
	globalCollection.openers = append(globalCollection.openers, opener)
	globalCollection.mu.Unlock()
	log.SetDebug(os.Stderr)
	//log.Debugf(" opener \"%s\" registered (%s)\n", opener.Name(), opener.Description())
	return nil
}

func Open(s string, lock bool) (io.ReadCloser, error) {
	log.Debugf(" openers: request to open %s\n", s)

	var eMax float32
	var oChosen Opener
	for _, o := range globalCollection.openers {
		e := o.Evaluate(s)
		//log.Debugf(" openers: evaluate %s with \"%s\": %v\n", s, o.Name(), e)
		if e > eMax {
			eMax = e
			oChosen = o
		}

	}
	if eMax == 0.0 {
		if !strings.Contains(s, "://") {
			return nil, fs.ErrNotExist
		} else {
			return nil, errors.New("No adequate opener found.")
		}
	}
	log.Debugf(" openers: chosen one is \"%s\"\n", oChosen.Name())
	return oChosen.Open(s, lock)
}

func ListOpeners() []string {
	var l []string
	for _, o := range globalCollection.openers {
		l = append(l, o.Name()+": "+o.Description())
	}
	return l
}
