package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/etsy/hound/api"
	"github.com/etsy/hound/config"
	"github.com/etsy/hound/searcher"
	"github.com/etsy/hound/ui"
)

const gracefulShutdownSignal = syscall.SIGTERM

var (
	info_log  *log.Logger
	error_log *log.Logger
)

func makeSearchers(cfg *config.Config) (map[string]*searcher.Searcher, bool, error) {
	// Ensure we have a dbpath
	if _, err := os.Stat(cfg.DbPath); err != nil {
		if err := os.MkdirAll(cfg.DbPath, os.ModePerm); err != nil {
			return nil, false, err
		}
	}

	searchers, errs, err := searcher.MakeAll(cfg)
	if err != nil {
		return nil, false, err
	}

	if len(errs) > 0 {
		// NOTE: This mutates the original config so the repos
		// are not even seen by other code paths.
		for name, _ := range errs {
			delete(cfg.Repos, name)
		}

		return searchers, false, nil
	}

	return searchers, true, nil
}

func handleShutdown(shutdownCh <-chan os.Signal, searchers map[string]*searcher.Searcher) {
	go func() {
		<-shutdownCh
		info_log.Printf("Graceful shutdown requested...")
		for _, s := range searchers {
			s.Stop()
		}

		for _, s := range searchers {
			s.Wait()
		}

		os.Exit(0)
	}()
}

func registerShutdownSignal() <-chan os.Signal {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, gracefulShutdownSignal)
	return shutdownCh
}

func makeTemplateData(cfg *config.Config) (interface{}, error) {
	var data struct {
		ReposAsJson string
	}

	res := map[string]*config.Repo{}
	for name, repo := range cfg.Repos {
		res[strings.ToLower(name)] = repo
	}

	b, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	data.ReposAsJson = string(b)
	return &data, nil
}

func runHttp(
	addr string,
	dev bool,
	cfg *config.Config,
	idx map[string]*searcher.Searcher) error {
	m := http.DefaultServeMux

	h, err := ui.Content(dev, cfg)
	if err != nil {
		return err
	}

	m.Handle("/", h)
	api.Setup(m, idx)
	return http.ListenAndServe(addr, m)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	info_log = log.New(os.Stdout, "", log.LstdFlags)
	error_log = log.New(os.Stderr, "", log.LstdFlags)

	flagConf := flag.String("conf", "config.json", "")
	flagAddr := flag.String("addr", ":6080", "")
	flagDev := flag.Bool("dev", false, "")

	flag.Parse()

	var cfg config.Config
	if err := cfg.LoadFromFile(*flagConf); err != nil {
		panic(err)
	}

	// It's not safe to be killed during makeSearchers, so register the
	// shutdown signal here and defer processing it until we are ready.
	shutdownCh := registerShutdownSignal()
	idx, ok, err := makeSearchers(&cfg)
	if err != nil {
		log.Panic(err)
	}
	if !ok {
		info_log.Println("Some repos failed to index, see output above")
	} else {
		info_log.Println("All indexes built!")
	}

	handleShutdown(shutdownCh, idx)

	host := *flagAddr
	if strings.HasPrefix(host, ":") {
		host = "localhost" + host
	}

	info_log.Printf("running server at http://%s...\n", host)

	if err := runHttp(*flagAddr, *flagDev, &cfg, idx); err != nil {
		panic(err)
	}
}
