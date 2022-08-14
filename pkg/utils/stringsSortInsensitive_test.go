package utils_test

import (
	"reflect"
	"testing"

	"github.com/batmac/ccat/pkg/utils"
)

func TestSortStringsCaseInsensitive(t *testing.T) {
	tests := []struct {
		name     string
		data     []string
		expected []string
	}{
		{"empty", []string{}, []string{}},
		{"empty string", []string{""}, []string{""}},
		{"string", []string{"abc", "abc-Z", "abc-abc"}, []string{"abc", "abc-abc", "abc-Z"}},
		{
			"digit",
			[]string{"8", "7", "2", "4", "5", "3", "1", "6", "0", "9"},
			[]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
		},
		{
			"case",
			[]string{"c-8", "C-7", "C-2", "c-4", "c-5", "c-3", "C-1", "C-6", "C-0", "c-9"},
			[]string{"C-0", "C-1", "C-2", "c-3", "c-4", "c-5", "C-6", "C-7", "c-8", "c-9"},
		},
		{"case-2", []string{"UPPER", "lower_9", "LOWER_1"}, []string{"LOWER_1", "lower_9", "UPPER"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils.SortStringsCaseInsensitive(tt.data)
			if !reflect.DeepEqual(tt.data, tt.expected) {
				t.Fatalf("%v failed: got %#v, expected %#v", tt.name, tt.data, tt.expected)
			}
		})
	}
}
