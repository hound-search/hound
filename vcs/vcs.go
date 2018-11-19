package vcs

import (
	"fmt"
	"log"
	"os"

	"github.com/etsy/hound/config"
)

// A collection that maps vcs names to their underlying
// factory. A factory allows the vcs to have unserialized
// json config passed in to be parsed.
var drivers = make(map[string]func(c []byte) (Driver, error))

// A "plugin" for each vcs that supports the very limited set of vcs
// operations that hound needs.
type Driver interface {

	// Clone a new working directory.
	Clone(dir, url string) (string, error)

	// Pull new changes from the server and update the working directory.
	Pull(dir string, ref string) (string, error)

	// Return the revision at the head of the vcs directory.
	HeadRev(dir string) (string, error)

	// Return a list of special filenames that should not be indexed.
	SpecialFiles() []string
}

// An API to interact with a vcs working directory. This is
// what clients will interact with.
type WorkDir struct {
	Driver
}

// Register a new vcs driver under 1 or more names.
func Register(fn func(c []byte) (Driver, error), names ...string) {
	if fn == nil {
		log.Panic("vcs: cannot register nil factory")
	}

	for _, name := range names {
		drivers[name] = fn
	}
}

// Create a new WorkDir from the name and configuration data.
func New(name string, cfg []byte) (*WorkDir, error) {
	f := drivers[name]
	if f == nil {
		return nil, fmt.Errorf("vcs: %s is not a valid vcs driver.", name)
	}

	d, err := f(cfg)
	if err != nil {
		return nil, err
	}

	return &WorkDir{d}, nil
}

func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

// A utility method that carries out the common operation of cloning
// if the working directory is absent and pulling otherwise.
func (w *WorkDir) PullOrClone(dir string, repo *config.Repo) (string, error) {
	if exists(dir) {
		return w.Pull(dir, repo.Ref)
	}
	return w.Clone(dir, repo.Url)
}
