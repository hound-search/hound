package vcs

import (
	"fmt"
	"strings"

	"github.com/etsy/hound/config"
)

func init() {
	Register(newLocal, "local")
}

type LocalDriver struct{}

func newLocal(b []byte) (Driver, error) {
	return &LocalDriver{}, nil
}

func (g *LocalDriver) WorkingDirForRepo(dbpath string, repo *config.Repo) (string, error) {
	return strings.TrimPrefix(repo.Url, "file://"), nil
}

func (g *LocalDriver) HeadRev(dir string) (string, error) {
	return "n/a", nil
}

func (g *LocalDriver) Pull(dir string) (string, error) {
	return g.HeadRev(dir)
}

func (g *LocalDriver) Clone(dir, url string) (string, error) {
	// For local driver Clone() is only called when the directory
	// pointed by url is not found.
	err := fmt.Errorf("Location %s not found.", url)
	fmt.Print(err)
	return "", err
}

func (g *LocalDriver) SpecialFiles() []string {
	return []string{
		".bzr",
		".git",
		".hg",
		".svn",
	}
}
