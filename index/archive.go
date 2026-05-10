package index

import (
	"context"
	"io/fs"
	"os"

	"github.com/hound-search/hound/codesearch/index"
	"github.com/hound-search/hound/config"
	"github.com/mholt/archiver/v4"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func isArchiveRepo(repo *config.Repo) bool {
	return stringInSlice(repo.Vcs, []string{"archive"})
}

func openArchive(repo *config.Repo) (fs.FS, error) {
	ctx := context.Background()
	return archiver.FileSystem(ctx, repo.Url)
}

func indexArchive(opt *IndexOptions, repo *config.Repo, ix *index.IndexWriter) ([]*ExcludedFile, error) {
	fsys, err := openArchive(repo)
	if err != nil {
		return nil, err
	}

	excluded := []*ExcludedFile{}

	err = fs.WalkDir(fsys, ".", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		name := info.Name()
		// already relative paths in archive
		rel := path

		// Is this file considered "special", this means it's not even a part
		// of the source repository (like .git or .svn).
		if containsString(opt.SpecialFiles, name) {
			if info.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if path == "." {
			// special case for archives
			return nil
		}

		if opt.ExcludeDotFiles && name[0] == '.' {
			if info.IsDir() {
				return fs.SkipDir
			}

			excluded = append(excluded, &ExcludedFile{
				rel,
				reasonDotFile,
			})
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if info.Type()&os.ModeType & ^os.ModeSymlink != 0 {
			excluded = append(excluded, &ExcludedFile{
				rel,
				reasonInvalidMode,
			})
			return nil
		}

		// is text file
		{
			r, err := fsys.Open(path)
			if err != nil {
				return err
			}
			defer r.Close()

			txt, err := isTextReader(r)
			if err != nil {
				return err
			}

			if !txt {
				excluded = append(excluded, &ExcludedFile{
					rel,
					reasonNotText,
				})
				return nil
			}
		}

		r, err := fsys.Open(path)
		if err != nil {
			return err
		}
		defer r.Close()

		reasonForExclusion := ix.Add(rel, r)
		if reasonForExclusion != "" {
			excluded = append(excluded, &ExcludedFile{rel, reasonForExclusion})
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return excluded, nil
}
