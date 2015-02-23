package vcs

import (
	"fmt"
	"os"
)

type Driver interface {
	Clone(dir, url, branch string) (string, error)
	Pull(dir, branch string) (string, error)
	HeadHash(dir string) (string, error)
}

var drivers = make(map[string]Driver)

func RegisterVCS(name string, driver Driver) {
	if driver == nil {
		panic("vcs: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("vcs: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func Clone(vcs, dir, url, branch string) (string, error) {
	driver, ok := drivers[vcs]
	if !ok {
		return "", fmt.Errorf("vcs: unknown driver %q", vcs)
	}

	return driver.Clone(dir, url, branch)
}

func Pull(vcs, dir, branch string) (string, error) {
	driver, ok := drivers[vcs]
	if !ok {
		return "", fmt.Errorf("vcs: unknown driver %q", vcs)
	}

	return driver.Pull(dir, branch)
}

func HeadHash(vcs, dir string) (string, error) {
	driver, ok := drivers[vcs]
	if !ok {
		return "", fmt.Errorf("vcs: unknown driver %q", vcs)
	}
	return driver.HeadHash(dir)
}

func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func PullOrClone(vcs, dir, url, branch string) (string, error) {
	if exists(dir) {
		return Pull(vcs, dir, branch)
	}

	return Clone(vcs, dir, url, branch)
}
