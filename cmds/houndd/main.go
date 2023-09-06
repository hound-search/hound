package main

import (
	"flag"
	"fmt"
	"github.com/blang/semver/v4"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"

	"github.com/hound-search/hound/config"
	"github.com/hound-search/hound/searcher"
	"github.com/hound-search/hound/web"
)

var (
	info_log   *log.Logger
	error_log  *log.Logger
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func makeSearchers(cfg *config.Config, searchers map[string]*searcher.Searcher) (bool, error) {
	// Ensure we have a dbpath
	if _, err := os.Stat(cfg.DbPath); err != nil {
		if err := os.MkdirAll(cfg.DbPath, os.ModePerm); err != nil {
			return false, err
		}
	}

	errs, err := searcher.MakeAll(cfg, searchers)
	if err != nil {
		return false, err
	}

	if len(errs) > 0 {
		// NOTE: This mutates the original config so the repos
		// are not even seen by other code paths.
		for name := range errs {
			delete(cfg.Repos, name)
		}

		return false, nil
	}

	return true, nil
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

func registerShutdownSignal() chan os.Signal {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGTERM)
	signal.Notify(shutdownCh, syscall.SIGINT)
	return shutdownCh
}

func unregisterShutdownSignal(shutdownCh chan os.Signal) {
	if shutdownCh == nil {
		return
	}
	signal.Stop(shutdownCh)
	shutdownCh = nil
}

// TODO: Automatically increment this when building a release
func getVersion() semver.Version {
	return semver.Version{
		Major: 0,
		Minor: 7,
		Patch: 1,
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	info_log = log.New(os.Stdout, "", log.LstdFlags)
	error_log = log.New(os.Stderr, "", log.LstdFlags)

	flagConf := flag.String("conf", "config.json", "")
	flagAddr := flag.String("addr", ":6080", "")
	flagDev := flag.Bool("dev", false, "")
	flagVer := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *flagVer {
		fmt.Printf("houndd v%s", getVersion())
		os.Exit(0)
	}

	idx := make(map[string]*searcher.Searcher)

	var cfg config.Config

	var shutdownCh chan os.Signal
	shutdownCh = nil
	var configUpdateLock sync.Mutex

	loadConfig := func(server *web.Server) {
		configUpdateLock.Lock()
		defer configUpdateLock.Unlock()
		// store existing cfg to check if it's changed
		cfgJson, _ := cfg.ToJsonString()

		if err := cfg.LoadFromFile(*flagConf); err != nil {
			panic(err)
		}
		// unregister shutdown signal to create a new one
		unregisterShutdownSignal(shutdownCh)

		// It's not safe to be killed during makeSearchers, so register the
		// shutdown signal here and defer processing it until we are ready.
		shutdownCh = registerShutdownSignal()

		ok, err := makeSearchers(&cfg, idx)
		if err != nil {
			log.Panic(err)
		}
		if !ok {
			info_log.Println("Some repos failed to index, see output above")
		} else {
			info_log.Println("All indexes built!")
		}
		// handle shutdown signal now
		handleShutdown(shutdownCh, idx)
		if server != nil {
			// if server has been passed, now check if cfg json has been changed
			newCfgJson, err := cfg.ToJsonString()
			if (err == nil) && (cfgJson != newCfgJson) {
				// cfg json changed, update the server
				info_log.Println("configJson updated, reloading server")
				if err := server.UpdateServeWithIndex(idx); err != nil {
					error_log.Printf("updating server failed for some reason: %s", err)
				}
			}
		}
	}

	// Start the web server on a background routine.
	ws := web.Start(&cfg, *flagAddr, *flagDev)

	host := *flagAddr
	if strings.HasPrefix(host, ":") { //nolint
		host = "localhost" + host
	}
	info_log.Printf("started server without indexes at http://%s\n", host)

	// Initial config load
	loadConfig(nil)
	info_log.Printf("loaded config")

	if *flagDev {
		info_log.Printf("[DEV] starting webpack-dev-server at localhost:8080...")
		webpack := exec.Command("./node_modules/.bin/webpack-dev-server", "--mode", "development")
		webpack.Dir = basepath + "/../../"
		webpack.Stdout = os.Stdout
		webpack.Stderr = os.Stderr

		if err := webpack.Start(); err != nil {
			error_log.Println(err)
		}
	}

	// watch for config file changes
	configWatcher := config.NewWatcher(*flagConf)
	configWatcher.OnChange(
		func(fsnotify.Event) {
			loadConfig(ws)
		},
	)

	info_log.Printf("running server at http://%s\n", host)

	// Fully enable the web server now that we have indexes
	panic(ws.ServeWithIndex(idx))
}
