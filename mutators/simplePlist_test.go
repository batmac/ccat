package mutators_test

import (
	"testing"

	"github.com/batmac/ccat/mutators"
	"github.com/batmac/ccat/utils"
)

func Test_simplePlist(t *testing.T) {
	var tests = []struct {
		name, decoded, encoded string
	}{
		{"test",
			`<?xml version="1.0" encoding="UTF-8"?>
		 <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
		 <plist version="1.0">
		 <dict>
			 <key>testKey</key>
			 <string>testString</string>
			 <key>Label</key>
			 <string>testLabel</string>
			 <key>testDict</key>
			 <dict>
				 <key>testKeyInDict</key>
				 <true/>
			 </dict>
			 <key>testArray</key>
			 <array>
				 <string>testStringInArray</string>
			 </array>
		 </dict>
		 </plist>`,
			"Label: testLabel\ntestArray:\n- testStringInArray\ntestDict:\n  testKeyInDict: true\ntestKey: testString"},
	}

	f := "plist2Y"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); utils.DeleteSpaces(got) != utils.DeleteSpaces(tt.encoded) {
				t.Errorf("%s = %v, want %v", f, got, tt.encoded)
			}
		})
	}
}
