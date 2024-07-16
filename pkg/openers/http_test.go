//go:build !fileonly
// +build !fileonly

package openers_test

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/batmac/ccat/pkg/globalctx"
	"github.com/batmac/ccat/pkg/openers"
)

var (
	hostTest         = "127.0.0.1"
	portTest         int
	insecurePortTest int
	m                sync.Mutex
)

func init() {
	// rand.Seed(time.Now().UnixNano())

	go func() {
		// find an available port
		for i := 0; i < 100; i++ {
			m.Lock()
			portTest = 10000 + rand.Intn(55000) //#nosec G404
			m.Unlock()
			_ = http.ListenAndServe(hostTest+":"+strconv.Itoa(portTest), SimpleHandler()) //nolint:gosec
		}
	}()
	go func() {
		// find an available port
		for i := 0; i < 100; i++ {
			m.Lock()
			insecurePortTest = 10000 + rand.Intn(55000) //#nosec G404
			m.Unlock()
			err := http.ListenAndServeTLS(hostTest+":"+strconv.Itoa(insecurePortTest), "testdata/server.crt", "testdata/server.bin", SimpleHandler()) //nolint:gosec
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
	// give the routines some time to find available ports.
	time.Sleep(1 * time.Second)
}

func Test_http(t *testing.T) {
	m.Lock()
	type args struct {
		s    string
		lock bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"ok", args{"http://" + hostTest + ":" + strconv.Itoa(portTest), false}, false},
		{"fake", args{"http://fakefakefake", false}, true},
	}
	m.Unlock()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := openers.Open(tt.args.s, tt.args.lock)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			// exhaust the reader
			_, _ = io.Copy(io.Discard, f)
			f.Close()
		})
	}
}

func Test_https(t *testing.T) {
	m.Lock()
	type args struct {
		s    string
		lock bool
	}
	tests := []struct {
		name     string
		args     args
		insecure bool
		wantErr  bool
	}{
		{"ko", args{s: "https://" + hostTest + ":" + strconv.Itoa(insecurePortTest), lock: false}, false, true},
		{"ok", args{s: "https://" + hostTest + ":" + strconv.Itoa(insecurePortTest), lock: false}, true, false},
		{"fake", args{"https://fakefakefake", false}, false, true},
		{"fake", args{"https://fakefakefake", false}, true, true},
	}
	m.Unlock()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			globalctx.Set("insecure", tt.insecure)
			f, err := openers.Open(tt.args.s, tt.args.lock)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			// exhaust the reader
			_, _ = io.Copy(io.Discard, f)
			f.Close()
		})
	}
}

func Simple(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "200 hello", http.StatusOK)
}

func SimpleHandler() http.Handler {
	return http.HandlerFunc(Simple)
}
