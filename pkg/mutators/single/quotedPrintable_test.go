package mutators_test

import (
	"testing"

	"github.com/batmac/ccat/pkg/mutators"
)

func Test_simpleUnQP(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"hello", "hello", "hello"},
		{"empty", "", ""},
		{
			"wikipedia",
			"J'interdis aux marchands de vanter trop leurs marchandises. Car ils se font=\n" +
				" vite p=C3=A9dagogues et t'enseignent comme but ce qui n'est par essence qu=\n" +
				"'un moyen, et te trompant ainsi sur la route =C3=A0 suivre les voil=C3=\n" +
				"=A0 bient=C3=B4t qui te d=C3=A9gradent, car si leur musique est vulgaire il=\n" +
				"s te fabriquent pour te la vendre une =C3=A2me vulgaire.\n" +
				"=E2=80=94=E2=80=89Antoine de Saint-Exup=C3=A9ry, Citadelle (1948)",
			"J'interdis aux marchands de vanter trop leurs marchandises. Car ils se font " +
				"vite pédagogues et t'enseignent comme but ce qui n'est par essence qu'un moyen," +
				" et te trompant ainsi sur la route à suivre les voilà bientôt qui te dégradent," +
				" car si leur musique est vulgaire ils te fabriquent pour te la vendre une âme vulgaire.\n" +
				"— Antoine de Saint-Exupéry, Citadelle (1948)",
		},
	}
	f := "unqp"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); got != tt.encoded {
				t.Errorf("%s = %v, want %v", f, got, tt.encoded)
			}
		})
	}
}

func Test_simpleQP(t *testing.T) {
	tests := []struct {
		name, decoded, encoded string
	}{
		{"hello", "hello", "hello"},
		{"empty", "", ""},
		{
			"wikipedia",
			"J'interdis aux marchands de vanter trop leurs marchandises. Car ils se font " +
				"vite pédagogues et t'enseignent comme but ce qui n'est par essence qu'un moyen," +
				" et te trompant ainsi sur la route à suivre les voilà bientôt qui te dégradent," +
				" car si leur musique est vulgaire ils te fabriquent pour te la vendre une âme vulgaire.\n" +
				"— Antoine de Saint-Exupéry, Citadelle (1948)",
			"J'interdis aux marchands de vanter trop leurs marchandises. Car ils se font=\r\n" +
				" vite p=C3=A9dagogues et t'enseignent comme but ce qui n'est par essence qu=\r\n" +
				"'un moyen, et te trompant ainsi sur la route =C3=A0 suivre les voil=C3=A0 b=\r\n" +
				"ient=C3=B4t qui te d=C3=A9gradent, car si leur musique est vulgaire ils te =\r\n" +
				"fabriquent pour te la vendre une =C3=A2me vulgaire.\r\n" +
				"=E2=80=94=E2=80=89Antoine de Saint-Exup=C3=A9ry, Citadelle (1948)",
		},
	}
	f := "qp"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mutators.Run(f, tt.decoded); got != tt.encoded {
				t.Errorf("%s = %v, want %v", f, got, tt.encoded)
			}
		})
	}
}
