package config

import (
	"encoding/json"
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

	// Ensure that global VCS config vals are merged
	repo := cfg.Repos["SomeGitRepo"]
	vcsConfigBytes := repo.VcsConfig()
	var vcsConfigVals map[string]interface{}
	json.Unmarshal(vcsConfigBytes, &vcsConfigVals)  //nolint
	if detectRef, ok := vcsConfigVals["detect-ref"]; !ok || !detectRef.(bool) {
		t.Error("global detectRef vcs config setting not set for repo")
	}

	repo = cfg.Repos["GitRepoWithDetectRefDisabled"]
	vcsConfigBytes = repo.VcsConfig()
	json.Unmarshal(vcsConfigBytes, &vcsConfigVals)  //nolint
	if detectRef, ok := vcsConfigVals["detect-ref"]; !ok || detectRef.(bool) {
		t.Error("global detectRef vcs config setting not overriden by repo-level setting")
	}

}
