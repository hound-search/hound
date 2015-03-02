package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	defaultMsBetweenPoll = 30000
	defaultVcs           = "git"
	defaultBaseUrl       = "{url}/blob/master/{path}{anchor}"
	defaultAnchor        = "#L{line}"
)

type UrlPattern struct {
	BaseUrl string `json:"base-url"`
	Anchor  string `json:"anchor"`
}

type Repo struct {
	Url            string      `json:"url"`
	MsBetweenPolls int         `json:"ms-between-poll"`
	Vcs            string      `json:"vcs"`
	UrlPattern     *UrlPattern `json:"url-pattern"`
}

type Config struct {
	DbPath string           `json:"dbpath"`
	Repos  map[string]*Repo `json:"repos"`
}

// Populate missing config values with default values.
func initRepo(r *Repo) {
	if r.MsBetweenPolls == 0 {
		r.MsBetweenPolls = defaultMsBetweenPoll
	}

	if r.Vcs == "" {
		r.Vcs = defaultVcs
	}

	if r.UrlPattern == nil {
		r.UrlPattern = &UrlPattern{
			BaseUrl: defaultBaseUrl,
			Anchor:  defaultAnchor,
		}
	} else {
		if r.UrlPattern.BaseUrl == "" {
			r.UrlPattern.BaseUrl = defaultBaseUrl
		}

		if r.UrlPattern.Anchor == "" {
			r.UrlPattern.Anchor = defaultAnchor
		}
	}
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
		initRepo(repo)
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
