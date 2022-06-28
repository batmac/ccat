package utils_test

import (
	"testing"

	"github.com/batmac/ccat/utils"
)

func TestRemoveScheme(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{""}, ""},
		{"empty://", args{"empty://"}, ""},
		{"http", args{"http://localhost"}, "localhost"},
		{":port", args{"http://:80"}, ":80"},
		{"http:port", args{"http://localhost:80"}, "localhost:80"},
		{"k8s", args{"https://hello.world.svc.cluster.local:9000"}, "hello.world.svc.cluster.local:9000"},
		{"mc", args{"mc://localhost"}, "localhost"},
		{"tcp", args{"tcp://localhost"}, "localhost"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.RemoveScheme(tt.args.s); got != tt.want {
				t.Errorf("RemoveScheme() = %v, want %v", got, tt.want)
			}
		})
	}
}
