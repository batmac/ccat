package mutators_test

import (
	"strings"
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

func Test_simpleSha256(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"empty", "", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"simple", "{}", "44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a"},
		{"indented", `{"hi":"hi"}`, "14ab3e46196e88fbe28ab26935f700edf6860e780a110657dd6e607ffe1bb630"},
		{"indented2", `{"hi": 1}`, "14cbe38a4dc585f1c558cb57488790caaf394bf2ed7138d58c6bb5370e413948"},
		{"indented3", "   { \n \"hi\" :    1 \n    }", "eb0d3743537f40eef62fa4cbc53a4c929fc0f23d1d056eb8d1eecdc84e00af28"},
		{"indented4", "  \n\n { \n\n \n  \n \"hi\"  \n: \n   1 \n}", "d841b4fdecbe48dc59ed72c1a7572002916d8026b9891bbfa9d1e4fe331f78a1"},
	}

	f := "sha256"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); strings.TrimSuffix(got, "\n") != tt.encoded {
				t.Errorf("%s = %v, want %v", f, strings.TrimSuffix(got, "\n"), tt.encoded)
			}
		})
	}
}

func Test_simpleSha256std(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"empty", "", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"simple", "{}", "44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a"},
		{"indented", `{"hi":"hi"}`, "14ab3e46196e88fbe28ab26935f700edf6860e780a110657dd6e607ffe1bb630"},
		{"indented2", `{"hi": 1}`, "14cbe38a4dc585f1c558cb57488790caaf394bf2ed7138d58c6bb5370e413948"},
		{"indented3", "   { \n \"hi\" :    1 \n    }", "eb0d3743537f40eef62fa4cbc53a4c929fc0f23d1d056eb8d1eecdc84e00af28"},
		{"indented4", "  \n\n { \n\n \n  \n \"hi\"  \n: \n   1 \n}", "d841b4fdecbe48dc59ed72c1a7572002916d8026b9891bbfa9d1e4fe331f78a1"},
	}

	f := "sha256std"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); strings.TrimSuffix(got, "\n") != tt.encoded {
				t.Errorf("%s = %v, want %v", f, strings.TrimSuffix(got, "\n"), tt.encoded)
			}
		})
	}
}

func Test_simpleXxhash(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"empty", "", "ef46db3751d8e999"},
		{"simple", "{}", "2e1472b57af294d1"},
		{"indented", `{"hi":"hi"}`, "932ee0befae3b64b"},
		{"indented2", `{"hi": 1}`, "81f26dc0993fd1d2"},
		{"indented3", "   { \n \"hi\" :    1 \n    }", "c6f9d585fdf2863b"},
		{"indented4", "  \n\n { \n\n \n  \n \"hi\"  \n: \n   1 \n}", "f165adb16aa69b92"},
	}

	f := "xxhash"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); strings.TrimSuffix(got, "\n") != tt.encoded {
				t.Errorf("%s = %v, want %v", f, strings.TrimSuffix(got, "\n"), tt.encoded)
			}
		})
	}
}
