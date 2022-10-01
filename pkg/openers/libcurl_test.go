//go:build cgo && libcurl && !fileonly
// +build cgo,libcurl,!fileonly

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
	curlHostTest         = "127.0.0.1"
	curlPortTest         int
	curlInsecurePortTest int

	curlm sync.Mutex
)

func init() {
	go func() {
		rand.Seed(time.Now().UnixNano())
		// find an available port
		for {
			curlm.Lock()
			curlPortTest = 10000 + rand.Intn(55000) //#nosec G404
			curlm.Unlock()
			_ = http.ListenAndServe(curlHostTest+":"+strconv.Itoa(curlPortTest), curlSimpleHandler())
		}
	}()
	go func() {
		// find an available port
		for i := 0; i < 100; i++ {
			m.Lock()
			curlInsecurePortTest = 10000 + rand.Intn(55000) //#nosec G404
			m.Unlock()
			err := http.ListenAndServeTLS(curlHostTest+":"+strconv.Itoa(curlInsecurePortTest), "testdata/server.crt", "testdata/server.bin", curlSimpleHandler())
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
	// give the routine some time to find an available port.
	time.Sleep(1 * time.Second)
}

func Test_libcurlHttp(t *testing.T) {
	curlm.Lock()
	type args struct {
		s    string
		lock bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"ok", args{"curlhttp://" + curlHostTest + ":" + strconv.Itoa(curlPortTest), false}, false},
		{"fake", args{"curlhttp://fakefakefake", false}, true},
	}
	curlm.Unlock()
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

func Test_libcurlHttps(t *testing.T) {
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
		{"ko", args{s: "curlhttps://" + curlHostTest + ":" + strconv.Itoa(curlInsecurePortTest), lock: false}, false, true},
		{"ok", args{s: "curlhttps://" + curlHostTest + ":" + strconv.Itoa(curlInsecurePortTest), lock: false}, true, false},
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

func curlSimple(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "200 hello", http.StatusOK)
}

func curlSimpleHandler() http.Handler {
	return http.HandlerFunc(curlSimple)
}
