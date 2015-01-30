package vcs

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"log"
)

func init() {
	RegisterVCS("svn", &SvnDriver{})
}

type SvnDriver struct{}

func pull(dir string) error {
	cmd := exec.Command("svn", "update")
	cmd.Dir = dir
	return cmd.Run()
}

func (g *SvnDriver) HeadHash(dir string) (string, error) {
	cmd := exec.Command(
		"svn",
		"log",
		"--limit 1")
	cmd.Dir = dir
	cmd.Stdout = ioutil.Discard
	//TODO create a better hash using the actual revision
	//TODO do not ignore out, but it's always returning exit status 1 - normal for SVN?
	cmd.Run()

	return "HEAD", nil
}

func (g *SvnDriver) Pull(dir string) (string, error) {
	if err := pull(dir); err != nil {
		return "", err
	}

	return g.HeadHash(dir)
}

func (g *SvnDriver) Clone(dir, url string) (string, error) {
	par, dirName := filepath.Split(dir)
	log.Printf("Checkout into %s from %s...\n", dir, url)
	cmd := exec.Command(
		"svn",
		"checkout",
		url,
		dirName)
	cmd.Dir = par
	cmd.Stdout = ioutil.Discard
	//TODO do not ignore out, but it's always returning exit status 1 - normal for SVN?
	cmd.Run()

	return g.HeadHash(dir)
}

