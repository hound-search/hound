package git

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func pull(dir string) error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = dir
	return cmd.Run()
}

func HeadHash(dir string) (string, error) {
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

func Pull(dir string) (string, error) {
	if err := pull(dir); err != nil {
		return "", err
	}

	return HeadHash(dir)
}

func Clone(dir, url string) (string, error) {
	par, rep := filepath.Split(dir)
	cmd := exec.Command(
		"git",
		"clone",
		url,
		rep)
	cmd.Dir = par
	cmd.Stdout = ioutil.Discard
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return HeadHash(dir)
}

func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func PullOrClone(dir, url string) (string, error) {
	if exists(dir) {
		return Pull(dir)
	}
	return Clone(dir, url)
}
