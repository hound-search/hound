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

func parentDir() string {
	return filepath.Dir(thisDir())
}

func buildIndex(url, rev string, opt *IndexOptions) (*IndexRef, error) {
	dir, err := ioutil.TempDir(os.TempDir(), "hound")
	if err != nil {
		return nil, err
	}

	return Build(opt, dir, parentDir(), url, rev)
}

func TestSearch(t *testing.T) {
	var opt IndexOptions

	// Build an index
	ref, err := buildIndex(url, rev, &opt)
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
	var opt IndexOptions

	ref, err := buildIndex(url, rev, &opt)
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
	var opt IndexOptions

	ref, err := buildIndex(url, rev, &opt)
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
		t.Fatal("expected rev of %s, got %s", rev, r.Rev)
	}

	idx, err := r.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer idx.Close()
}

func TestExclude(t *testing.T) {
	opt := &IndexOptions {
		ExcludePatterns: []string{"index"},
	}

	// Build an index
	ref, err := buildIndex(url, rev, opt)
	if err != nil {
		t.Fatal(err)
	}
	defer ref.Remove()

	// Make sure the ref can be opened.
	idx, err := ref.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer idx.Close()

	// search for this test method. should not be found as it's excluded
	res, err := idx.Search("TestExclude", &SearchOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if res.FilesWithMatch > 0 {
		t.Fatalf("This file was excluded, but indexed : %v", res)
	}
}

func TestIncludeThis(t *testing.T) {
	opt := &IndexOptions {
		IncludePatterns: []string{"index"},
	}

	// Build an index
	ref, err := buildIndex(url, rev, opt)
	if err != nil {
		t.Fatal(err)
	}
	defer ref.Remove()

	// Make sure the ref can be opened.
	idx, err := ref.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer idx.Close()

	// search for this test method.
	res, err := idx.Search("TestIncludeThis", &SearchOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if res.FilesWithMatch != 1 {
		t.Fatalf("Failed to find this test method: %v", res)
	}
}

func TestIncludeOther(t *testing.T) {
	opt := &IndexOptions {
		IncludePatterns: []string{"config"},
	}

	// Build an index
	ref, err := buildIndex(url, rev, opt)
	if err != nil {
		t.Fatal(err)
	}
	defer ref.Remove()

	// Make sure the ref can be opened.
	idx, err := ref.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer idx.Close()

	// search for this test method.
	res, err := idx.Search("TestIncludeOther", &SearchOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if res.FilesWithMatch > 0 {
		t.Fatalf("This file was not included, but indexed: %v", res)
	}
}

func TestExcludeAndInclude(t *testing.T) {
	opt := &IndexOptions {
		ExcludePatterns: []string{"index"},
		IncludePatterns: []string{"index"},
	}

	// Build an index
	ref, err := buildIndex(url, rev, opt)
	if err != nil {
		t.Fatal(err)
	}
	defer ref.Remove()

	// Make sure the ref can be opened.
	idx, err := ref.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer idx.Close()

	// search for this test method.
	res, err := idx.Search("TestExcludeAndInclude", &SearchOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if res.FilesWithMatch > 0 {
		t.Fatalf("This file should be excluded, but was indexed: %v", res)
	}
}
