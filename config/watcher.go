package config

import (
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// WatcherListenerFunc defines the signature for listner functions
type WatcherListenerFunc func(fsnotify.Event)

// Watcher watches for configuration updates and provides hooks for
// triggering post events
type Watcher struct {
	listeners []WatcherListenerFunc
}

// NewWatcher returns a new file watcher
func NewWatcher(cfgPath string) *Watcher {
	log.Printf("setting up watcher for %s", cfgPath)
	w := Watcher{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Panic(err)
		}
		defer watcher.Close()
		// Event listener setup
		eventWG := sync.WaitGroup{}
		eventWG.Add(1)
		go func() {
			defer eventWG.Done()
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						// events channel is closed
						log.Printf("error: events channel is closed\n")
						return
					}
					// only trigger on creates and writes of the watched config file
					if event.Name == cfgPath && event.Op&fsnotify.Write == fsnotify.Write {
						log.Printf("change in config file (%s) detected\n", cfgPath)
						for _, listener := range w.listeners {
							listener(event)
						}
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						// errors channel is closed
						log.Printf("error: errors channel is closed\n")
						return
					}
					log.Println("error:", err)
					return
				}
			}
		}()
		// add config file
		if err := watcher.Add(cfgPath); err != nil {
			log.Fatalf("failed to watch %s", cfgPath)
		}
		// setup is complete
		wg.Done()
		// wait for the event listener to complete before exiting
		eventWG.Wait()
	}()
	// wait for watcher setup to complete
	wg.Wait()
	return &w
}

// OnChange registers a listener function to be called if a file changes
func (w *Watcher) OnChange(listener WatcherListenerFunc) {
	w.listeners = append(w.listeners, listener)
}
