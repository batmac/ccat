package pipeline_test

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/batmac/ccat/pkg/mutators/pipeline"
	_ "github.com/batmac/ccat/pkg/mutators/simple"

	"github.com/batmac/ccat/pkg/utils"
)

func TestNewPipeline(t *testing.T) {
	type args struct {
		description string
		out         string
		in          string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"dummy", args{"dummy", "in", "in"}, false},
		{"dummy,dummy", args{"dummy,dummy", "in 2", "in 2"}, false},
		{"dummy,dummy,dummy", args{"dummy,dummy,dummy", "in 3", "in 3"}, false},
		{"dummy,dummy,dummy,dummy", args{"dummy,dummy,dummy,dummy", "in 4", "in 4"}, false},
		{"dummy,dummy,dummy,dummy,dummy", args{"dummy,dummy,dummy,dummy,dummy", "in 5", "in 5"}, false},
		{"dummy,dummy,dummy,dummy,dummy,dummy", args{"dummy,dummy,dummy,dummy,dummy,dummy", "in 6", "in 6"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := utils.NopStringWriteCloser{Builder: strings.Builder{}}
			err := pipeline.NewPipeline(tt.args.description, &n, io.NopCloser(strings.NewReader(tt.args.in)))
			time.Sleep(1 * time.Second)
			if pipeline.Wait(); (err != nil) != tt.wantErr || tt.args.out != n.String() {
				t.Errorf("NewPipeline() error = %v, wantErr %v - out = %v/, want %v", err, tt.wantErr, n.String(), tt.args.out)
			}
		})
	}
}
