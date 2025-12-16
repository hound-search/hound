package vcs

import (
	"io"
	"io/fs"
	"path"
	"slices"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type gitFilesystem struct {
	repo *gogit.Repository
	root *object.Tree
}

type gitFileinfo struct {
	raw *object.File
}

type gitDirinfo struct {
	name string
}

func NewGitFilesystem(dir, ref string) (FileSystem, error) {
	repo, err := gogit.PlainOpen(dir)
	if err != nil {
		return nil, err
	}

	rev, err := repo.ResolveRevision(plumbing.Revision(ref))
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(*rev)
	if err != nil {
		return nil, err
	}

	root, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	return &gitFilesystem{
		repo: repo,
		root: root,
	}, nil
}

func (fs *gitFilesystem) Open(name string) (io.ReadCloser, error) {
	file, err := fs.root.File(name)
	if err != nil {
		return nil, err
	}
	return file.Reader()
}

func (fs *gitFilesystem) Walk(fn FileSystemWalkFunc) error {
	seenDirs := make(map[string]bool)

	return fs.root.Files().ForEach(func(f *object.File) error {
		n := f.Name
		var createDirs []string
		if f.Mode != filemode.Dir {
			n = path.Dir(n)
		}
		for n != "" && !seenDirs[n] {
			seenDirs[n] = true
			createDirs = append(createDirs, n)
			n = path.Dir(n)
		}
		slices.Reverse(createDirs)
		for _, createDir := range createDirs {
			if err := fn(createDir, &gitDirinfo{n}, nil); err != nil {
				return err
			}
		}

		if f.Mode != filemode.Dir {
			return fn(f.Name, &gitFileinfo{f}, nil)
		} else {
			return nil
		}
	})
}

func (fi *gitFileinfo) Name() string {
	return path.Base(fi.raw.Name)
}

func (fi *gitFileinfo) IsDir() bool {
	return fi.Mode().IsDir()
}

func (fi *gitFileinfo) Mode() fs.FileMode {
	mode, _ := fi.raw.Mode.ToOSFileMode()
	return mode
}

func (di *gitDirinfo) Name() string {
	return path.Base(di.name)
}
func (di *gitDirinfo) IsDir() bool {
	return di.Mode().IsDir()
}
func (di *gitDirinfo) Mode() fs.FileMode {
	return fs.FileMode(0o755 | fs.ModeDir)
}
