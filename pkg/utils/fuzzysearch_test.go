package utils_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/utils"
)

func TestFuzzySearch(t *testing.T) {
	type args struct {
		str       string
		strList   []string
		threshold float32
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test", args{"tsest", []string{"first", "test", "notest", "cat", ""}, 0.5}, "test", false},
		{"tet", args{"tet", []string{"first", "test", "notest", "cat", ""}, 0.5}, "test", false},
		{"estt", args{"estt", []string{"first", "test", "notest", "cat", ""}, 0.5}, "test", false},
		{"markdown", args{"markdown", []string{"test", "md", "notest", "cat", ""}, 0.5}, "md", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.FuzzySearch(tt.args.str, tt.args.strList, tt.args.threshold)
			if (err != nil) != tt.wantErr {
				t.Errorf("FuzzySearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FuzzySearch() = %v, want %v", got, tt.want)
			}
		})
	}
}
