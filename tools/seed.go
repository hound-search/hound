package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/etsy/hound/config"
	"github.com/google/go-github/github"
	"github.com/tj/docopt"
	"golang.org/x/oauth2"
)

const (
	usage = `Seeds an initial config.json file, with all the repositories
	owned by the given organizations.

Usage:
  seed --token=<token> <organizations>... [--dbpath=<dbpath>] [--indexers=<i>]
  seed -h | --help
  seed --version

Options:
  --dbpath=<dbpath>       Database path [default: data].
  --indexers=<i>          Max concurrent indexers [default: 2].
  --token=<token>         Github oauth token.
  -h --help               Show this screen.
  --version               Show version.`

	version = "0.1.0"
)

func main() {
	arguments, err := docopt.Parse(usage, nil, true, version, false)
	check(err)

	config := readConfig(arguments)

	f, err := os.Create("config.json")
	check(err)
	defer f.Close()

	b, err := json.Marshal(config)
	check(err)

	// Prettifying is not required, but makes it easier for humans to read.
	var buf bytes.Buffer
	err = json.Indent(&buf, b, "  ", "\t")
	check(err)
	_, err = buf.WriteTo(f)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func readConfig(arguments map[string]interface{}) *config.Config {
	token := arguments["--token"].(string)
	organizations := arguments["<organizations>"].([]string)
	dbpath := arguments["--dbpath"].(string)
	indexers, err := strconv.Atoi(arguments["--indexers"].(string))
	check(err)

	repos := make(map[string]*config.Repo)

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	opt := &github.RepositoryListByOrgOptions{
		Type:        "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for _, org := range organizations {
		for {
			newRepos, resp, err := client.Repositories.ListByOrg(org, opt)
			check(err)

			for _, newRepo := range newRepos {
				url := fmt.Sprintf("git@github.com:%s/%s.git", org, *newRepo.Name)
				repos[*newRepo.Name] = &config.Repo{
					Url: url,
				}
			}

			if resp.NextPage == 0 {
				break
			}
			opt.ListOptions.Page = resp.NextPage
		}
	}

	return &config.Config{
		DbPath: dbpath,
		Repos:  repos,
		MaxConcurrentIndexers: indexers,
	}
}
