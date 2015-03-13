package vcs

import (
	"testing"
)

// Tests that the svn driver is able to parse its config.
func TestSvnConfig(t *testing.T) {
	cfg := `{"username" : "svn_username", "password" : "svn_password"}`

	d, err := New("svn", []byte(cfg))
	if err != nil {
		t.Fatal(err)
	}

	svn := d.Driver.(*SVNDriver)
	if svn.Username != "svn_username" {
		t.Fatalf("expected username of \"svn_username\", got %s", svn.Username)
	}

	if svn.Password != "svn_password" {
		t.Fatalf("expected password of \"svn_password\", got %s", svn.Password)
	}
}
