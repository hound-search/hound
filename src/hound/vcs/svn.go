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
	Username string
	Password string
}

func newSvn(b []byte) Driver {
	var d SVNDriver

	if err := json.Unmarshal(b, &d); err != nil {
		// TODO(knorton): I guess these really need to reutrn error
		log.Panic(err)
	}

	return &d
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

func (g *SVNDriver) Pull(dir string) (string, error) {
	cmd := exec.Command(
		"svn",
		"update",
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
