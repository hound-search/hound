package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	defaultMsBetweenPoll = 30000
	defaultVcs           = "git"
)

type Repo struct {
	Url            string   `json:"url"`
	Branches       []string `json:"branches"`
	MsBetweenPolls int      `json:"ms-between-poll"`
	Vcs            string   `json:"vcs"`
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
