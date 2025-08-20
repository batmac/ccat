package mutators_test

import (
	"strings"
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

func Test_simpleSha3_224(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"empty", "", "6b4e03423667dbb73b6e15454f0eb1abd4597f9a1b078e3f5b5a6bc7"},
		{"abc", "abc", "e642824c3f8cf24ad09234ee7d3c766fc9a3a5168d0c94ad73b46fdf"},
		{"simple", "{}", "9661d47169c2d8014928d7c850283d9e74790f614155a853595b36fa"},
		{"alphanumeric", "abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq", "8a24108b154ada21c9fd5574494479ba5c7e7ab76ef264ead0fcce33"},
		{"million_a", strings.Repeat("a", 1000000), "d69335b93325192e516a912e6d19a15cb51c6ed5c15243e7a7fd653c"},
	}

	f := "sha3-224"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); strings.TrimSuffix(got, "\n") != tt.encoded {
				t.Errorf("%s = %v, want %v", f, strings.TrimSuffix(got, "\n"), tt.encoded)
			}
		})
	}
}

func Test_simpleSha3_256(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"empty", "", "a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a"},
		{"abc", "abc", "3a985da74fe225b2045c172d6bd390bd855f086e3e9d525b46bfe24511431532"},
		{"simple", "{}", "840eb7aa2a9935de63366bacbe9d97e978a859e93dc792a0334de60ed52f8e99"},
		{"alphanumeric", "abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq", "41c0dba2a9d6240849100376a8235e2c82e1b9998a999e21db32dd97496d3376"},
		{"million_a", strings.Repeat("a", 1000000), "5c8875ae474a3634ba4fd55ec85bffd661f32aca75c6d699d0cdcb6c115891c1"},
	}

	f := "sha3-256"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); strings.TrimSuffix(got, "\n") != tt.encoded {
				t.Errorf("%s = %v, want %v", f, strings.TrimSuffix(got, "\n"), tt.encoded)
			}
		})
	}
}

func Test_simpleSha3_384(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"empty", "", "0c63a75b845e4f7d01107d852e4c2485c51a50aaaa94fc61995e71bbee983a2ac3713831264adb47fb6bd1e058d5f004"},
		{"abc", "abc", "ec01498288516fc926459f58e2c6ad8df9b473cb0fc08c2596da7cf0e49be4b298d88cea927ac7f539f1edf228376d25"},
		{"simple", "{}", "3763179c317870e535852680e55663659ac85034d5f99bbf31642fcfeaf2754a76350a19db6a8111a83609be2f1901ca"},
		{"alphanumeric", "abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq", "991c665755eb3a4b6bbdfb75c78a492e8c56a22c5c4d7e429bfdbc32b9d4ad5aa04a1f076e62fea19eef51acd0657c22"},
		{"million_a", strings.Repeat("a", 1000000), "eee9e24d78c1855337983451df97c8ad9eedf256c6334f8e948d252d5e0e76847aa0774ddb90a842190d2c558b4b8340"},
	}

	f := "sha3-384"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); strings.TrimSuffix(got, "\n") != tt.encoded {
				t.Errorf("%s = %v, want %v", f, strings.TrimSuffix(got, "\n"), tt.encoded)
			}
		})
	}
}

func Test_simpleSha3_512(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"empty", "", "a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26"},
		{"abc", "abc", "b751850b1a57168a5693cd924b6b096e08f621827444f70d884f5d0240d2712e10e116e9192af3c91a7ec57647e3934057340b4cf408d5a56592f8274eec53f0"},
		{"simple", "{}", "c1802e6b9670927ebfddb7f67b3824642237361f07db35526c42c555ffd2dbe74156c366e1550ef8c0508a6cc796409a7194a59bba4d300a6182b483d315a862"},
		{"alphanumeric", "abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq", "04a371e84ecfb5b8b77cb48610fca8182dd457ce6f326a0fd3d7ec2f1e91636dee691fbe0c985302ba1b0d8dc78c086346b533b49c030d99a27daf1139d6e75e"},
		{"million_a", strings.Repeat("a", 1000000), "3c3a876da14034ab60627c077bb98f7e120a2a5370212dffb3385a18d4f38859ed311d0a9d5141ce9cc5c66ee689b266a8aa18ace8282a0e0db596c90b0a7b87"},
	}

	f := "sha3-512"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); strings.TrimSuffix(got, "\n") != tt.encoded {
				t.Errorf("%s = %v, want %v", f, strings.TrimSuffix(got, "\n"), tt.encoded)
			}
		})
	}
}