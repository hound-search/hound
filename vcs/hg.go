package vcs

import (
	"bytes"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
)

func init() {
	Register(newHg, "hg", "mercurial")
}

type MercurialDriver struct{}

func newHg(b []byte) (Driver, error) {
	return &MercurialDriver{}, nil
}

func (g *MercurialDriver) HeadRev(dir string) (string, error) {
	cmd := exec.Command(
		"hg",
		"log",
		"-r",
		".",
		"--template",
		"{node}")
	cmd.Dir = dir
	r, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	defer r.Close()

	if err := cmd.Start(); err != nil {
		return "", err
	}

	var buf bytes.Buffer

	if _, err := io.Copy(&buf, r); err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), cmd.Wait()
}

func (g *MercurialDriver) Pull(dir string) (string, error) {
	cmd := exec.Command("hg", "pull", "-u")
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return g.HeadRev(dir)
}

func (g *MercurialDriver) Clone(dir, url string) (string, error) {
	par, rep := filepath.Split(dir)
	cmd := exec.Command(
		"hg",
		"clone",
		url,
		rep)
	cmd.Dir = par
	cmd.Stdout = ioutil.Discard
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return g.HeadRev(dir)
}

func (g *MercurialDriver) SpecialFiles() []string {
	return []string{
		".hg",
	}
}
