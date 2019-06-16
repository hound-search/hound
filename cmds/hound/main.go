package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/user"
	"regexp"

	"github.com/hound-search/hound/client"
	"github.com/hound-search/hound/index"
)

// A uninitialized variable that can be defined during the build process with
// -ldflags -X main.defaultHouse addr. This should remain uninitialized.
var defaultHost string

// a convenience method for creating a new presenter that is either
// ack-like or grep-like.
func newPresenter(likeGrep bool) client.Presenter {
	if likeGrep {
		return client.NewGrepPresenter(os.Stdout)
	}

	return client.NewAckPresenter(os.Stdout)
}

// the paths we will attempt to load config from
var configPaths = []string{
	"/etc/hound.conf",
	"$HOME/.hound",
}

// Attempt to populate a client.Config from the json found in
// filename.
func loadConfigFrom(filename string, cfg *client.Config) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	return json.NewDecoder(r).Decode(cfg)
}

// Attempt to populate a client.Config from the json found in
// any of the configPaths.
func loadConfig(cfg *client.Config) error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	env := map[string]string{
		"HOME": u.HomeDir,
	}

	for _, path := range configPaths {
		err = loadConfigFrom(os.Expand(path, func(name string) string {
			return env[name]
		}), cfg)

		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			return err
		}
	}

	return nil
}

// A simple way to determine what the default value should be
// for the --host flag.
func defaultFlagForHost() string {
	if defaultHost != "" {
		return defaultHost
	}
	return "localhost:6080"
}

func main() {
	flagHost := flag.String("host", defaultFlagForHost(), "")
	flagRepos := flag.String("repos", "*", "")
	flagFiles := flag.String("files", "", "")
	flagContext := flag.Int("context", 2, "")
	flagCase := flag.Bool("ignore-case", false, "")
	flagStats := flag.Bool("show-stats", false, "")
	flagGrep := flag.Bool("like-grep", false, "")

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		return
	}

	pat := index.GetRegexpPattern(flag.Arg(0), *flagCase)

	reg, err := regexp.Compile(pat)
	if err != nil {
		// TODO(knorton): Better error reporting
		log.Panic(err)
	}

	cfg := client.Config{
		Host:        *flagHost,
		HttpHeaders: nil,
	}

	if err := loadConfig(&cfg); err != nil {
		log.Panic(err)
	}

	res, repos, err := client.SearchAndLoadRepos(&cfg,
		flag.Arg(0),
		*flagRepos,
		*flagFiles,
		*flagContext,
		*flagCase,
		*flagStats)
	if err != nil {
		log.Panic(err)
	}

	if err := newPresenter(*flagGrep).Present(reg, *flagContext, repos, res); err != nil {
		log.Panic(err)
	}
}
