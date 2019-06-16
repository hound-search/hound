package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/hound-search/hound/config"
	"github.com/hound-search/hound/index"
)

type Response struct {
	Results map[string]*index.SearchResponse
	Stats   *struct {
		FilesOpened int
		Duration    int
	} `json:",omitempty"`
}

type Presenter interface {
	Present(
		re *regexp.Regexp,
		ctx int,
		repos map[string]*config.Repo,
		res *Response) error
}

type Config struct {
	HttpHeaders map[string]string `json:"http-headers"`
	Host        string            `json:"host"`
}

// Extract a repo name from the given url.
func repoNameFromUrl(uri string) string {
	ax := strings.LastIndex(uri, "/")
	if ax < 0 {
		return ""
	}

	name := uri[ax+1:]
	if strings.HasSuffix(name, ".git") {
		name = name[:len(name)-4]
	}

	bx := strings.LastIndex(uri[:ax-1], "/")
	if bx < 0 {
		return name
	}

	return fmt.Sprintf("%s/%s", uri[bx+1:ax], name)
}

// Find the proper name for the given repo using the map of repo
// information.
func repoNameFor(repos map[string]*config.Repo, repo string) string {
	data := repos[repo]
	if data == nil {
		return repo
	}

	name := repoNameFromUrl(data.Url)
	if name == "" {
		return repo
	}

	return name
}

func doHttpGet(cfg *Config, uri string) (*http.Response, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	for key, val := range cfg.HttpHeaders {
		if strings.ToLower(key) == "host" {
			req.Host = val
		} else {
			req.Header.Set(key, val)
		}
	}

	var c http.Client
	return c.Do(req)
}

// Executes a search on the API running on host.
func Search(r *Response, cfg *Config, pattern, repos, files string, context int, ignoreCase, stats bool) error {
	u := fmt.Sprintf("http://%s/api/v1/search?%s",
		cfg.Host,
		url.Values{
			"q":     {pattern},
			"repos": {repos},
			"files": {files},
			"ctx":   {fmt.Sprintf("%d", context)},
			"i":     {fmt.Sprintf("%t", ignoreCase)},
			"stats": {fmt.Sprintf("%t", stats)},
		}.Encode())

	res, err := doHttpGet(cfg, u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Status %d", res.StatusCode)
	}

	return json.NewDecoder(res.Body).Decode(r)
}

// Load the list of repositories from the API running on host.
func LoadRepos(repos map[string]*config.Repo, cfg *Config) error {
	res, err := doHttpGet(cfg, fmt.Sprintf("http://%s/api/v1/repos", cfg.Host))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(&repos)
}

// Execute a search and load the list of repositories in parallel on the host.
func SearchAndLoadRepos(cfg *Config, pattern, repos, files string, context int, ignoreCase, stats bool) (*Response, map[string]*config.Repo, error) {
	chs := make(chan error)
	var res Response
	go func() {
		chs <- Search(&res, cfg, pattern, repos, files, context, ignoreCase, stats)
	}()

	chr := make(chan error)
	rep := map[string]*config.Repo{}
	go func() {
		chr <- LoadRepos(rep, cfg)
	}()

	// must ready both channels before returning to avoid routine/channel leak.
	errS, errR := <-chs, <-chr

	if errS != nil {
		return nil, nil, errS
	}

	if errR != nil {
		return nil, nil, errR
	}

	return &res, rep, nil
}
