package index

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

const (
	url = "https://www.etsy.com/"
	rev = "r420"
)

func thisDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}

func buildIndex(url, rev string) (*IndexRef, error) {
	dir, err := ioutil.TempDir(os.TempDir(), "hound")
	if err != nil {
		return nil, err
	}

	var opt IndexOptions

	return Build(&opt, dir, thisDir(), url, rev)
}

func TestSearch(t *testing.T) {
	// Build an index
	ref, err := buildIndex(url, rev)
	if err != nil {
		t.Fatal(err)
	}
	defer ref.Remove()

	// Make sure the metadata in the ref is good.
	if ref.Rev != rev {
		t.Fatalf("expected rev of %s, got %s", rev, ref.Rev)
	}

	if ref.Url != url {
		t.Fatalf("expected url of %s got %s", url, ref.Url)
	}

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

func TestRemove(t *testing.T) {
	ref, err := buildIndex(url, rev)
	if err != nil {
		t.Fatal(err)
	}

	if err := ref.Remove(); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(ref.Dir()); err == nil {
		t.Fatalf("Remove did not delete directory: %s", ref.Dir())
	}
}

func TestRead(t *testing.T) {
	ref, err := buildIndex(url, rev)
	if err != nil {
		t.Fatal(err)
	}
	defer ref.Remove()

	r, err := Read(ref.Dir())
	if err != nil {
		t.Fatal(err)
	}

	if r.Url != url {
		t.Fatalf("expected url of %s, got %s", url, r.Url)
	}

	if r.Rev != rev {
		t.Fatalf("expected rev of %s, got %s", rev, r.Rev)
	}

	idx, err := r.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer idx.Close()
}
