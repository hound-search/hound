package vcs

import (
	"os"
	"testing"
)

// TODO(knorton): Write tests for the vcs interactions

// Just make sure all drivers are tolerant of nil
func TestNilConfigs(t *testing.T) {
	for name, _ := range drivers { //nolint
		d, err := New(name, nil)
		if err != nil {
			t.Fatal(err)
		}

		if d == nil {
			t.Fatalf("vcs: %s returned a nil driver", name)
		}
	}
}

func TestIsWriteable(t *testing.T) {
	dir := t.TempDir()
	if writeable, err := IsWriteable(dir); !writeable {
		t.Fatalf("%s is not writeable but should be: %s", dir, err)
	}
	if err := os.Chmod(dir, 0444); err != nil {
		t.Fatal(err)
	}

	if writeable, err := IsWriteable(dir); writeable {
		t.Fatalf("%s is writeable but should not be: %s", dir, err)
	}
}
