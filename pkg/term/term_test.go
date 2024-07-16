package term_test

import (
	"os"
	"testing"

	"github.com/batmac/ccat/pkg/term"
)

func TestIsStdoutTerminal(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := term.IsStdoutTerminal(); got != tt.want {
				t.Errorf("IsStdoutTerminal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsStdinTerminal(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := term.IsStdinTerminal(); got != tt.want {
				t.Errorf("IsStdinTerminal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTerminalSize(t *testing.T) {
	tests := []struct {
		name       string
		wantWidth  int
		wantHeight int
		wantErr    bool
	}{
		{"", 80, 24, false}, // default
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWidth, gotHeight, err := term.GetTerminalSize()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTerminalSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWidth != tt.wantWidth {
				t.Errorf("GetTerminalSize() gotWidth = %v, want %v", gotWidth, tt.wantWidth)
			}
			if gotHeight != tt.wantHeight {
				t.Errorf("GetTerminalSize() gotHeight = %v, want %v", gotHeight, tt.wantHeight)
			}
		})
	}
}

func TestClearScreen(t *testing.T) {
	t.Run("donotpanicplease", func(_ *testing.T) {
		term.ClearScreen()
	})
}

func TestSupportedColors(t *testing.T) {
	tests := []struct {
		name string
		want uint
	}{
		{"", 8}, // default
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := term.SupportedColors(); got < tt.want {
				t.Errorf("SupportedColors() = %v, want >= %v", got, tt.want)
			}
		})
	}
}

func TestReadLine(t *testing.T) {
	w, cleaner := redirectStdin(t)
	defer cleaner()

	testStr := "test"
	t.Run("ReadLine", func(t *testing.T) {
		_, _ = w.WriteString(testStr + "\n")
		if got, err := term.ReadLine("input:"); got != testStr || err != nil {
			t.Errorf("ReadLine() = %v, err = %v want %v", got, err, testStr)
		}
	})
}

// we can't test ReadPassword() because it needs a real terminal to work (run some ioctl)

func redirectStdin(t *testing.T) (*os.File, func()) {
	t.Helper()
	r, w, _ := os.Pipe()
	saved := os.Stdin
	os.Stdin = r

	return w, func() {
		r.Close()
		w.Close()
		os.Stdin = saved
	}
}
