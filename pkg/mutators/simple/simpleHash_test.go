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

func Test_simpleXxh3(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"empty", "", "2d06800538d394c2"},
		{"simple", "{}", "1349cde127705c16"},
		{"indented", `{"hi":"hi"}`, "de1bfffbced9a6ae"},
		{"indented2", `{"hi": 1}`, "373643e61692f145"},
		{"indented3", "   { \n \"hi\" :    1 \n    }", "9548609af4b6cc0e"},
		{"indented4", "  \n\n { \n\n \n  \n \"hi\"  \n: \n   1 \n}", "660b7b829ae5faa9"},
	}

	f := "xxh3"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); strings.TrimSuffix(got, "\n") != tt.encoded {
				t.Errorf("%s = %v, want %v", f, strings.TrimSuffix(got, "\n"), tt.encoded)
			}
		})
	}
}

func Test_simpleSha1(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"empty", "", "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
		{"simple", "{}", "bf21a9e8fbc5a3846fb05b4fa0859e0917b2202f"},
		{"indented", `{"hi":"hi"}`, "119a4c40ea410512fc2d3ea886d2ab7767e82f59"},
		{"indented2", `{"hi": 1}`, "502ee1f19995b0994e4d5dadf36423bf6ae88350"},
		{"indented3", "   { \n \"hi\" :    1 \n    }", "6c228cebf277f662c786d213cf89ea464afe72b2"},
		{"indented4", "  \n\n { \n\n \n  \n \"hi\"  \n: \n   1 \n}", "0644b3d443e933b3cf7c2c6846334679f6c8079a"},
	}

	f := "sha1"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); strings.TrimSuffix(got, "\n") != tt.encoded {
				t.Errorf("%s = %v, want %v", f, strings.TrimSuffix(got, "\n"), tt.encoded)
			}
		})
	}
}

func Test_simpleMd5(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"empty", "", "d41d8cd98f00b204e9800998ecf8427e"},
		{"simple", "{}", "99914b932bd37a50b983c5e7c90ae93b"},
		{"indented", `{"hi":"hi"}`, "096ef7138af380ee63b7999f35c42ba4"},
		{"indented2", `{"hi": 1}`, "388c90b997822327b63892d0a934ecef"},
		{"indented3", "   { \n \"hi\" :    1 \n    }", "ca64ddcb03632577b0336345eac146b9"},
		{"indented4", "  \n\n { \n\n \n  \n \"hi\"  \n: \n   1 \n}", "bd7b6cb1227dcd70beaddbddbabcf6e3"},
	}

	f := "md5"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); strings.TrimSuffix(got, "\n") != tt.encoded {
				t.Errorf("%s = %v, want %v", f, strings.TrimSuffix(got, "\n"), tt.encoded)
			}
		})
	}
}
