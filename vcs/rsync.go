package vcs

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
)

func init() {
	Register(newRsync, "rsync")
}

type RsyncDriver struct{}

func newRsync(b []byte) (Driver, error) {
	return &RsyncDriver{}, nil
}

func (g *RsyncDriver) HeadRev(dir string) (string, error) {
	return "n/a", nil
}

func (g *RsyncDriver) Pull(dir string) (string, error) {
	return g.HeadRev(dir)
}

func (g *RsyncDriver) Clone(dir, url string) (string, error) {
	par, rep := filepath.Split(dir)
	cmd := exec.Command(
		"rsync",
		"-r",
		url[7:],
		rep)
	cmd.Dir = par
	cmd.Stdout = ioutil.Discard
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return g.HeadRev(dir)
}

func (g *RsyncDriver) SpecialFiles() []string {
	return []string{}
}
