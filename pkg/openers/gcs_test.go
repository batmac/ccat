package openers_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/openers"
)

const (
	// random gs:// URLs
	gsURL  = "gs://cloud-samples-data/bigquery/us-states/us-states-by-date.csv"
	gs2URL = "gs://tuf-root-staging/root.json"
)

func Test_gcs(t *testing.T) {
	type args struct {
		s    string
		lock bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "gs OK", args: args{gsURL, false}, wantErr: false},
		{name: "gs KO", args: args{gsURL + "fake", false}, wantErr: true},
		{name: "gs2 OK", args: args{gs2URL, false}, wantErr: false},
		{name: "gs2 KO", args: args{gs2URL + "fake", false}, wantErr: true},
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
