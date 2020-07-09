package vcs

import (
	"crypto/sha1"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func init() {
	Register(newFake, "fake")
}

type FakeDriver struct{}

func newFake(b []byte) (Driver, error) {
	return &FakeDriver{}, nil
}

func (g *FakeDriver) HeadRev(dir string) (string, error) {
	hsh := sha1.New()
	cmd := exec.Command(
		"tar",
		"cf", "-",
		"./",
	)
	cmd.Dir = dir
	cmd.Stdout = hsh
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v: %+v", cmd.Args, err)
	}

	return fmt.Sprintf("%x", hsh.Sum(nil)), nil
}

const fakeSpec = ".source"

func (g *FakeDriver) Pull(dir string) (string, error) {
	return g.HeadRev(dir)
}

func (g *FakeDriver) Clone(dir, url string) (string, error) {
	src := strings.TrimPrefix(url, "file://")
	if err := os.Symlink(src, dir); err != nil {
		return "", err
	}
	return g.Pull(dir)
}

func (g *FakeDriver) SpecialFiles() []string {
	return nil
}
