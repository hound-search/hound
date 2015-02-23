package vcs

import (
	"bytes"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

func init() {
	RegisterVCS("git", &GitDriver{})
}

type GitDriver struct{}

func (g *GitDriver) HeadHash(dir string) (string, error) {
	cmd := exec.Command(
		"git",
		"rev-parse",
		"HEAD")
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

func (g *GitDriver) Checkout(dir, branch string) error {
	cmd := exec.Command("git", "checkout", branch)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to git checkout %s, see output below\n%sContinuing...", dir, out)
		return err
	}
	return nil
}

func (g *GitDriver) Pull(dir, branch string) (string, error) {
	if branch == "" {
		branch = "master"
	}

	if err := g.Checkout(dir, branch); err != nil {
		return "", err
	}

	cmd := exec.Command("git", "pull")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to git pull %s, see output below\n%sContinuing...", dir, out)
		return "", err
	}

	return g.HeadHash(dir)
}

func (g *GitDriver) Clone(dir, url, branch string) (string, error) {
	if branch == "" {
		branch = "master"
	}
	par, rep := filepath.Split(dir)
	cmd := exec.Command(
		"git",
		"clone",
		"-b",
		branch,
		url,
		rep)
	cmd.Dir = par
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to clone %s, see output below\n%sContinuing...", url, out)
		return "", err
	}

	return g.HeadHash(dir)
}
