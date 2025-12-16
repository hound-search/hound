package vcs

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type dirFilesystem struct {
	dir string
}

func NewDirFilesystem(dir string) (FileSystem, error) {
	// Resolve the symbolic link
	if fi, err := os.Stat(dir); err == nil && fi.Mode()|os.ModeSymlink != 0 {
		if s, err := os.Readlink(dir); err == nil {
			dir = s
		}
	}

	return &dirFilesystem{dir: dir}, nil
}

func (dir *dirFilesystem) Open(name string) (io.ReadCloser, error) {
	if strings.HasPrefix(name, "/") {
		return nil, fmt.Errorf("Expected relative path, got absolute: %s", name)
	}
	return os.Open(path.Join(dir.dir, name))
}

func (dir *dirFilesystem) Walk(fn FileSystemWalkFunc) error {
	return filepath.Walk(dir.dir, func(path string, info fs.FileInfo, err error) error {
		rel, err := filepath.Rel(dir.dir, path)
		if err != nil {
			return err
		}

		return fn(rel, info, err)
	})
}
