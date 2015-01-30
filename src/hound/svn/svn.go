package svn

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"log"
)

func pull(dir string) error {
	cmd := exec.Command("svn", "update")
	cmd.Dir = dir
	return cmd.Run()
}

func HeadHash(dir string) (string, error) {
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

func Pull(dir string) (string, error) {
	if err := pull(dir); err != nil {
		return "", err
	}

	return HeadHash(dir)
}

func Clone(dir, url string) (string, error) {
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

	hash, err := HeadHash(dir)
	log.Printf("hash: %s err: %s\n", hash, err)
	return hash, err
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
