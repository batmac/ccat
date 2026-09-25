package selfupdate

import "testing"

// must match the -tags of the builds in .goreleaser.yaml
func Test_releaseTags(t *testing.T) {
	tests := []struct {
		goos, want string
	}{
		{"darwin", "plugins,keystore,libcurl_purego"},
		{"linux", "plugins,keystore"},
		{"windows", "plugins,keystore"},
	}
	for _, tt := range tests {
		if got := releaseTags(tt.goos); got != tt.want {
			t.Errorf("releaseTags(%q) = %q, want %q", tt.goos, got, tt.want)
		}
	}
}
