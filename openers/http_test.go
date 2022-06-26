//go:build !fileonly
// +build !fileonly

package openers_test

import (
	"net/http"
	"testing"

	"github.com/batmac/ccat/openers"
)

func Test_http(t *testing.T) {
	go func() { _ = http.ListenAndServe(":12344", SimpleHandler()) }()
	type args struct {
		s    string
		lock bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"exe", args{"http://" + "localhost:12344", false}, false},
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

func Simple(w http.ResponseWriter, r *http.Request) { http.Error(w, "200 hello", 200) }

func SimpleHandler() http.Handler { return http.HandlerFunc(Simple) }
