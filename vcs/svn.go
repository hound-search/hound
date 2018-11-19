package vcs

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

func init() {
	Register(newSvn, "svn", "subversion")
}

type SVNDriver struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func newSvn(b []byte) (Driver, error) {
	var d SVNDriver

	if b != nil {
		if err := json.Unmarshal(b, &d); err != nil {
			return nil, err
		}
	}

	return &d, nil
}

func (g *SVNDriver) HeadRev(dir string) (string, error) {
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

func (g *SVNDriver) Pull(dir string, _ string) (string, error) {
	cmd := exec.Command(
		"svn",
		"update",
		"--ignore-externals",
		"--username",
		g.Username,
		"--password",
		g.Password)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to SVN update %s, see output below\n%sContinuing...", dir, out)
		return "", err
	}

	return g.HeadRev(dir)
}

func (g *SVNDriver) Clone(dir, url string) (string, error) {
	par, rep := filepath.Split(dir)
	cmd := exec.Command(
		"svn",
		"checkout",
		"--ignore-externals",
		"--username",
		g.Username,
		"--password",
		g.Password,
		url,
		rep)
	cmd.Dir = par
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to checkout %s, see output below\n%sContinuing...", url, out)
		return "", err
	}

	return g.HeadRev(dir)
}

func (g *SVNDriver) SpecialFiles() []string {
	return []string{
		".svn",
	}
}
