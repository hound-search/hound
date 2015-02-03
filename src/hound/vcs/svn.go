package vcs

import (
	"bytes"
	"io"
	"strings"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"encoding/xml"
	"log"
)


type Info struct {
	Entry Entry `xml:"entry"`
}
type Entry struct {
	Revision string `xml:"revision,attr"`
	Url string `xml:"url"`
	RelativeUrl string `xml:"relative-url"`
	Repository Repository `xml:"repository"`
}
type Repository struct {
	Root string `xml:"root"`
	Uuid string `xml:"uuid"`
}

func init() {
	RegisterVCS("svn", &SvnDriver{})
}

type SvnDriver struct{}

func pull(dir string) error {
	cmd := exec.Command("svn", "update")
	cmd.Dir = dir
	return cmd.Run()
}

func (g *SvnDriver) HeadHash(dir string) (string, error) {
	//TODO is there a better way to extract revision number in golang?
	cmd := exec.Command(
		"svn", "info", "--xml")
	cmd.Dir = dir
	r, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	defer r.Close()

	if err := cmd.Start(); err != nil {
		return "", err
	}

	var buf bytes.Buffer

	if _, err := io.Copy(&buf, r); err != nil {
		return "", err
	}

	// TODO parse xml output
	info := Info{}
	if err := xml.Unmarshal(buf.Bytes(), &info); err != nil {
		log.Printf("error: %v", err)
		return "", err
	}

	hash := strings.TrimSpace(info.Entry.Repository.Uuid)
	log.Printf("hash: %s", hash)

	return hash, cmd.Wait()
}

func (g *SvnDriver) Pull(dir string) (string, error) {
	if err := pull(dir); err != nil {
		return "", err
	}

	return g.HeadHash(dir)
}

func (g *SvnDriver) Clone(dir, url string) (string, error) {
	par, dirName := filepath.Split(dir)
	cmd := exec.Command(
		"svn",
		"checkout",
		url,
		dirName)
	cmd.Dir = par
	cmd.Stdout = ioutil.Discard
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return g.HeadHash(dir)
}

