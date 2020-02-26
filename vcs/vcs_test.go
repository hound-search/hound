package vcs

import (
	"testing"
)

// TODO(knorton): Write tests for the vcs interactions

// Just make sure all drivers are tolerant of nil
func TestNilConfigs(t *testing.T) {
	for name := range drivers {
		d, err := New(name, nil)
		if err != nil {
			t.Fatal(err)
		}

		if d == nil {
			t.Fatalf("vcs: %s returned a nil driver", name)
		}
	}
}
