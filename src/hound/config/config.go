package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	defaultMsBetweenPoll = 30000
	defaultVcs           = "git"
	defaultUrlPattern    = "${url}/blob/master/${path}${anc}"
	defaultAnchorPattern = "#L${line}"
)

type Repo struct {
	Url            string `json:"url"`
	MsBetweenPolls int    `json:"ms-between-poll"`
	Vcs            string `json:"vcs"`
	UrlPattern     string `json:"urlpattern"`
	AnchorPattern  string `json:"anchorpattern"`
}

type Config struct {
	DbPath string           `json:"dbpath"`
	Repos  map[string]*Repo `json:"repos"`
}

func (c *Config) LoadFromFile(filename string) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := json.NewDecoder(r).Decode(c); err != nil {
		return err
	}

	if !filepath.IsAbs(c.DbPath) {
		path, err := filepath.Abs(
			filepath.Join(filepath.Dir(filename), c.DbPath))
		if err != nil {
			return err
		}
		c.DbPath = path
	}

	for _, repo := range c.Repos {
		if repo.MsBetweenPolls == 0 {
			repo.MsBetweenPolls = defaultMsBetweenPoll
		}
		if repo.Vcs == "" {
			repo.Vcs = defaultVcs
		}
		if repo.UrlPattern == "" {
			repo.UrlPattern = defaultUrlPattern
		}
		if repo.AnchorPattern == "" {
			repo.AnchorPattern = defaultAnchorPattern
		}
	}

	return nil
}

func (c *Config) ToJsonString() (string, error) {
	b, err := json.Marshal(c.Repos)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
