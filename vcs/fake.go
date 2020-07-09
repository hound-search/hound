package vcs

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func init() {
	Register(newFake, "fake")
}

type FakeDriver struct{}

func newFake(b []byte) (Driver, error) {
	return &FakeDriver{}, nil
}

func (g *FakeDriver) HeadRev(dir string) (string, error) {
	hsh := sha1.New()
	cmd := exec.Command(
		"tar",
		"cf", "-",
		".",
	)
	cmd.Dir = dir
	cmd.Stdout = hsh
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v: %+v", cmd.Args, err)
	}

	return fmt.Sprintf("%x", hsh.Sum(nil)), nil
}

const fakeSpec = ".source"

func (g *FakeDriver) Pull(dir string) (string, error) {
	if dir == "" || dir == "/" {
		return "", fmt.Errorf("empty dest dir %q", dir)
	}
	b, err := ioutil.ReadFile(filepath.Join(dir, fakeSpec))
	if err != nil {
		return "", err
	}
	src := string(b)
	if src == "" || src == "/" {
		return "", fmt.Errorf("empty source dir %q", src)
	}
	if err := os.MkdirAll(dir, 0750); err != nil {
		log.Printf("MkdirAll(%q): %+v", dir, err)
	}
	//log.Println("dir:", dir, "src:", src)
	if err := run("rsync fetch", dir,
		"rsync",
		"-a",
		"--delete-after",
		src+"/",
		"./",
	); err != nil {
		return "", err
	}
	if err := g.writeSource(dir, src); err != nil {
		return "", err
	}

	return g.HeadRev(dir)
}

func (g *FakeDriver) writeSource(dir, src string) error {
	return ioutil.WriteFile(filepath.Join(dir, fakeSpec), []byte(src), 0640)
}

func (g *FakeDriver) Clone(dir, url string) (string, error) {
	src := strings.TrimPrefix(url, "file://")
	os.MkdirAll(dir, 0750)
	if err := g.writeSource(dir, src); err != nil {
		return "", err
	}
	return g.Pull(dir)
}

func (g *FakeDriver) SpecialFiles() []string {
	return []string{fakeSpec}
}
