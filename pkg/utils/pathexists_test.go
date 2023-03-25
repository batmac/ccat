package utils_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/batmac/ccat/pkg/utils"
)

var executablePath string

func init() {
	var err error
	executablePath, err = os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func TestPathExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "root", args: args{path: "/"}, want: true},
		{name: "home", args: args{path: "~"}, want: true},
		{name: "executable", args: args{path: executablePath}, want: true},
		{name: "non-existent", args: args{path: "/non-existent"}, want: false},
		{name: "PWD", args: args{path: "$PWD"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.PathExists(tt.args.path); got != tt.want {
				t.Errorf("PathExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
