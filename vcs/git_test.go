package vcs

import "testing"

func TestGitConfigWithCustomRef(t *testing.T) {
	cfg := `{"ref": "custom"}`
	d, err := New("git", []byte(cfg))
	if err != nil {
		t.Fatal(err)
	}
	git := d.Driver.(*GitDriver)
	if git.Ref != "custom" {
		t.Fatalf("expected branch of \"custom\", got %s", git.Ref)
	}
}

func TestGitConfigWithoutRef(t *testing.T) {
	cfg := `{"option": "option"}`
	d, err := New("git", []byte(cfg))
	if err != nil {
		t.Fatal(err)
	}
	git := d.Driver.(*GitDriver)
	if git.Ref != "master" {
		t.Fatalf("expected branch of \"master\", got %s", git.Ref)
	}
}

func TestGitConfigWithoutAdditionalConfig(t *testing.T) {
	d, err := New("git", nil)
	if err != nil {
		t.Fatal(err)
	}
	git := d.Driver.(*GitDriver)
	if git.Ref != "master" {
		t.Fatalf("expected branch of \"master\", got %s", git.Ref)
	}
}
