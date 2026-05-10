package index

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hound-search/hound/config"
	"github.com/mholt/archiver/v4"
)

func buildIndexForZip() (*IndexRef, error) {
	rev := "n/a"

	dir, err := os.MkdirTemp(os.TempDir(), "hound")
	if err != nil {
		return nil, err
	}

	url := filepath.Join(dir, "archive.zip")
	if err := func() error {
		files, err := archiver.FilesFromDisk(nil, map[string]string{
			thisDir(): "",
		})

		out, err := os.Create(url)
		if err != nil {
			return err
		}
		defer out.Close()

		format := archiver.CompressedArchive{
			Archival: archiver.Zip{
				SelectiveCompression: true,
			},
		}

		return format.Archive(context.Background(), out, files)
	}(); err != nil {
		return nil, err
	}

	opt := &IndexOptions{
		SpecialFiles: []string{".git"},
	}

	return Build(opt, dir, "/not_existent_dir", &config.Repo{
		Url: url,
		Vcs: "archive",
	}, rev)
}

func TestSearchForZip(t *testing.T) {
	// Build an index
	ref, err := buildIndexForZip()
	if err != nil {
		t.Fatal(err)
	}
	defer ref.Remove() //nolint

	// Make sure the ref can be opened.
	idx, err := ref.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer idx.Close()

	// Make sure we can carry out a search
	if _, err := idx.Search("5a1c0dac2d9b3ea4085b30dd14375c18eab993d5", &SearchOptions{}); err != nil {
		t.Fatal(err)
	}
}
