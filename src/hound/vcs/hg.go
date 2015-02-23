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
	RegisterVCS("hg", &MercurialDriver{})
	RegisterVCS("mercurial", &MercurialDriver{})
}

type MercurialDriver struct{}

func (g *MercurialDriver) HeadHash(dir string) (string, error) {
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

func (g *MercurialDriver) Checkout(dir, branch string) error {
	cmd := exec.Command("hg", "checkout", branch)
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (g *MercurialDriver) Pull(dir, branch string) (string, error) {
	if branch == "" {
		branch = "default"
	}

	if err := g.Checkout(dir, branch); err != nil {
		return "", err
	}

	cmd := exec.Command("hg", "pull", "-u")
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return g.HeadHash(dir)
}

func (g *MercurialDriver) Clone(dir, url, branch string) (string, error) {
	if branch == "" {
		branch = "default"
	}
	par, rep := filepath.Split(dir)
	cmd := exec.Command(
		"hg",
		"clone",
		url,
		"-r",
		branch,
		rep)
	cmd.Dir = par
	cmd.Stdout = ioutil.Discard
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return g.HeadHash(dir)
}
