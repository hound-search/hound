package config

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hound-search/hound/vcs"
)

const exampleConfigFile = "config-example.json"

func rootDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..")
}

// Test that we can parse the example config file. This ensures that as we
// add examples, we don't muck them up.
func TestExampleConfigsAreValid(t *testing.T) {
	var cfg Config
	if err := cfg.LoadFromFile(filepath.Join(rootDir(), exampleConfigFile)); err != nil {
		t.Fatalf("Unable to parse %s: %s", exampleConfigFile, err)
	}

	if len(cfg.Repos) == 0 {
		t.Fatal("config has no repos")
	}

	// Ensure that each of the declared vcs's are legit
	for _, repo := range cfg.Repos {
		_, err := vcs.New(repo.Vcs, repo.VcsConfig())
		if err != nil {
			t.Fatal(err)
		}
	}
}
