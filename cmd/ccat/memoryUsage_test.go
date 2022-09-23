package main

import "testing"

func TestPrintMemUsage(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "donotpanicplease"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintMemUsage()
		})
	}
}
