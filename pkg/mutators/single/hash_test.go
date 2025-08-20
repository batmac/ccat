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

func Test_simpleXxh64(t *testing.T) {
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

	f := "xxh64"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); strings.TrimSuffix(got, "\n") != tt.encoded {
				t.Errorf("%s = %v, want %v", f, strings.TrimSuffix(got, "\n"), tt.encoded)
			}
		})
	}
}

func Test_simpleXxh32(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"empty", "", "02cc5d05"},
		{"simple", "{}", "18473489"},
		{"indented", `{"hi":"hi"}`, "980ceffc"},
		{"indented2", `{"hi": 1}`, "d1b0f9e8"},
		{"indented3", "   { \n \"hi\" :    1 \n    }", "e43b79fc"},
		{"indented4", "  \n\n { \n\n \n  \n \"hi\"  \n: \n   1 \n}", "8c1e37fa"},
	}

	f := "xxh32"
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

func Test_simpleSha224(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"empty", "", "d14a028c2a3a2bc9476102bb288234c415a2b01f828ea62ac5b3e42f"},
		{"abc", "abc", "23097d223405d8228642a477bda255b32aadbce4bda0b3f7e36c9da7"},
		{"simple", "{}", "5cdd15a873608087be07a41b7f1a04e96d3a66fe7a9b0faac71f8d05"},
		{"alphanumeric", "abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq", "75388b16512776cc5dba5da1fd890150b0c6455cb4f58b1952522525"},
		{"million_a", strings.Repeat("a", 1000000), "20794655980c91d8bbb4c1ea97618a4bf03f42581948b2ee4ee7ad67"},
	}

	f := "sha224"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); strings.TrimSuffix(got, "\n") != tt.encoded {
				t.Errorf("%s = %v, want %v", f, strings.TrimSuffix(got, "\n"), tt.encoded)
			}
		})
	}
}

func Test_simpleSha384(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"empty", "", "38b060a751ac96384cd9327eb1b1e36a21fdb71114be07434c0cc7bf63f6e1da274edebfe76f65fbd51ad2f14898b95b"},
		{"abc", "abc", "cb00753f45a35e8bb5a03d699ac65007272c32ab0eded1631a8b605a43ff5bed8086072ba1e7cc2358baeca134c825a7"},
		{"simple", "{}", "d2a23bc783e3aa38f401e13c7488505137c4954a7fd88331f1597c5ff71111dc807c7370a5b282c6da541c56ede69f30"},
		{"alphanumeric", "abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq", "3391fdddfc8dc7393707a65b1b4709397cf8b1d162af05abfe8f450de5f36bc6b0455a8520bc4e6f5fe95b1fe3c8452b"},
		{"million_a", strings.Repeat("a", 1000000), "9d0e1809716474cb086e834e310a4a1ced149e9c00f248527972cec5704c2a5b07b8b3dc38ecc4ebae97ddd87f3d8985"},
	}

	f := "sha384"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); strings.TrimSuffix(got, "\n") != tt.encoded {
				t.Errorf("%s = %v, want %v", f, strings.TrimSuffix(got, "\n"), tt.encoded)
			}
		})
	}
}

func Test_simpleSha512(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"empty", "", "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"},
		{"abc", "abc", "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a2192992a274fc1a836ba3c23a3feebbd454d4423643ce80e2a9ac94fa54ca49f"},
		{"simple", "{}", "27c74670adb75075fad058d5ceaf7b20c4e7786c83bae8a32f626f9782af34c9a33c2046ef60fd2a7878d378e29fec851806bbd9a67878f3a9f1cda4830763fd"},
		{"alphanumeric", "abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq", "204a8fc6dda82f0a0ced7beb8e08a41657c16ef468b228a8279be331a703c33596fd15c13b1b07f9aa1d3bea57789ca031ad85c7a71dd70354ec631238ca3445"},
		{"million_a", strings.Repeat("a", 1000000), "e718483d0ce769644e2e42c7bc15b4638e1f98b13b2044285632a803afa973ebde0ff244877ea60a4cb0432ce577c31beb009c5c2c49aa2e4eadb217ad8cc09b"},
	}

	f := "sha512"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); strings.TrimSuffix(got, "\n") != tt.encoded {
				t.Errorf("%s = %v, want %v", f, strings.TrimSuffix(got, "\n"), tt.encoded)
			}
		})
	}
}
