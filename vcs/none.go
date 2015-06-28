package vcs

import (
	"log"
	"os/exec"
	"path/filepath"
)

func init() {
	Register(newNoneVcs, "none")
}

type NoneVcsDriver struct{}

func newNoneVcs(b []byte) (Driver, error) {
	return &NoneVcsDriver{}, nil
}

func (g *NoneVcsDriver) HeadRev(dir string) (string, error) {
    return "n/a", nil
}

func (g *NoneVcsDriver) Pull(dir string) (string, error) {
	return g.HeadRev(dir)
}

func (g *NoneVcsDriver) Clone(dir, url string) (string, error) {
	par, rep := filepath.Split(dir)
	cmd := exec.Command(
		"cp",
		"-r",
		url[7:],
		rep)
	cmd.Dir = par
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to clone %s, see output below\n%sContinuing...", url, out)
		return "", err
	}

	return g.HeadRev(dir)
}

func (g *NoneVcsDriver) SpecialFiles() []string {
	return []string{}
}
