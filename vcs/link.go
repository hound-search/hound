package vcs

import (
	"log"
	"os"
	"strings"
)

func init() {
	Register(newLinkVcs, "link")
}

type LinkVcsDriver struct{}

func newLinkVcs(b []byte) (Driver, error) {
	return &LinkVcsDriver{}, nil
}

func (g *LinkVcsDriver) HeadRev(dir string) (string, error) {
	realdir, err := os.Readlink(dir)
	if err != nil {
		log.Printf("Failed to read symlink %s", dir)
		return "", err
	}

	stat, err := os.Stat(realdir)
	if err != nil {
		log.Printf("Failed to determine modification time of %s", realdir)
		return "", err
	}
	log.Printf("modtime %s", stat.ModTime().String())
	return stat.ModTime().String(), nil
}

func (g *LinkVcsDriver) Pull(dir string) (string, error) {
	return g.HeadRev(dir)
}

func (g *LinkVcsDriver) Clone(dir, url string) (string, error) {
	err := os.Symlink(strings.Replace(url, "file://", "", 1), dir)
	if err != nil {
		log.Printf("Failed to link %s, see output below\nContinuing...", url)
		return "", err
	}

	return g.HeadRev(dir)
}

func (g *LinkVcsDriver) SpecialFiles() []string {
	return []string{}
}
