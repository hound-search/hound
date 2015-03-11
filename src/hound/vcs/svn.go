package vcs

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func init() {
	RegisterVCS("svn", &SVNDriver{})
}

type SVNDriver struct{}

func (g *SVNDriver) HeadHash(dir string) (string, error) {
	cmd := exec.Command(
		"svnversion")
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

func (g *SVNDriver) Pull(dir string) (string, error) {
	cmd := exec.Command(
		"svn",
		"update",
		"--username",
		os.Getenv("SVN_USERNAME"),
		"--password",
		os.Getenv("SVN_PASSWORD"))
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to SVN update %s, see output below\n%sContinuing...", dir, out)
		return "", err
	}

	return g.HeadHash(dir)
}

func (g *SVNDriver) Clone(dir, url string) (string, error) {
	par, rep := filepath.Split(dir)
	cmd := exec.Command(
		"svn",
		"checkout",
		"--username",
		os.Getenv("SVN_USERNAME"),
		"--password",
		os.Getenv("SVN_PASSWORD"),
		url,
		rep)
	cmd.Dir = par
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to checkout %s, see output below\n%sContinuing...", url, out)
		return "", err
	}

	return g.HeadHash(dir)
}
