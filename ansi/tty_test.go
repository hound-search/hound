//go:build !windows

package ansi

import (
	"os"
	"testing"
)

func TestIsTTYTerminal(t *testing.T) {
	terminal, err := os.Open("/dev/tty")
	if err != nil {
		t.Skipf("no controlling terminal: %v", err)
	}
	defer func() { _ = terminal.Close() }()

	if !isTTY(terminal.Fd()) {
		t.Error("controlling terminal was not detected as a TTY")
	}
}

func TestIsTTYFile(t *testing.T) {
	file, err := os.CreateTemp(t.TempDir(), "tty-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = file.Close() }()

	if isTTY(file.Fd()) {
		t.Error("regular file was detected as a TTY")
	}
}
