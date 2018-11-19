package vcs

import (
	"crypto/sha1"
  "encoding/hex"
	"fmt"
  "path/filepath"
)

// Utility function for producing a hex encoded sha1 hash for a string.
func hashFor(name string) string {
	h := sha1.New()
	h.Write([]byte(name))
	return hex.EncodeToString(h.Sum(nil))
}

// Create a normalized name for the vcs directory of the repo.
func generateWorkingDir(dbpath string, url string) string {
	return filepath.Join(dbpath, fmt.Sprintf("vcs-%s", hashFor(url)))
}
