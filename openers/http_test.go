//go:build !fileonly
// +build !fileonly

package openers_test

import (
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/batmac/ccat/openers"
)

var (
	hostTest = "127.0.0.1"
	portTest int
)

func init() {
	go func() {
		rand.Seed(time.Now().UnixNano())
		// find an available port
		for {
			portTest = 10000 + rand.Intn(55000) //#nosec G404
			_ = http.ListenAndServe(hostTest+":"+strconv.Itoa(portTest), SimpleHandler())
		}
	}()
	// give the routine some time to find an available port.
	time.Sleep(3 * time.Second)
}

func Test_http(t *testing.T) {
	type args struct {
		s    string
		lock bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"exe", args{"http://" + hostTest + ":" + strconv.Itoa(portTest), false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := openers.Open(tt.args.s, tt.args.lock)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Simple(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "200 hello", 200)
}

func SimpleHandler() http.Handler {
	return http.HandlerFunc(Simple)
}
