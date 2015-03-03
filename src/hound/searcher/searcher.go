package searcher

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hound/config"
	"hound/index"
	"hound/vcs"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Searcher struct {
	idx  *index.Index
	lck  sync.RWMutex
	Repo *config.Repo
}

func (s *Searcher) swapIndexes(idx *index.Index) error {
	s.lck.Lock()
	defer s.lck.Unlock()

	oldIdx := s.idx
	s.idx = idx

	return oldIdx.Destroy()
}

func (s *Searcher) Search(pat string, opt *index.SearchOptions) (*index.SearchResponse, error) {
	s.lck.RLock()
	defer s.lck.RUnlock()
	return s.idx.Search(pat, opt)
}

func (s *Searcher) GetExcludedFiles() string {
	path := filepath.Join(s.idx.GetDir(), "excluded_files.json")
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Couldn't read excluded_files.json %v\n", err)
	}
	return string(dat)
}

func nextIndexName() string {
	r := uint64(rand.Uint32())<<32 | uint64(rand.Uint32())
	return fmt.Sprintf("idx-%08x", r)
}

func RemoveAllIndexes(dbpath string) error {
	dirs, err := filepath.Glob(filepath.Join(dbpath, "idx-*"))
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}

	return nil
}

func buildAndOpenIndex(dbpath, vcsDir string) (*index.Index, error) {
	idxDir := filepath.Join(dbpath, nextIndexName())
	if _, err := os.Stat(idxDir); err != nil {
		_, err := index.Build(idxDir, vcsDir)
		if err != nil {
			return nil, err
		}
	}

	return index.Open(idxDir)
}

func reportOnMemory() {
	var ms runtime.MemStats

	// Print out interesting heap info.
	runtime.ReadMemStats(&ms)
	fmt.Printf("HeapInUse = %0.2f\n", float64(ms.HeapInuse)/1e6)
	fmt.Printf("HeapIdle  = %0.2f\n", float64(ms.HeapIdle)/1e6)
}

func hashFor(name string) string {
	h := sha1.New()
	h.Write([]byte(name))
	return hex.EncodeToString(h.Sum(nil))
}

func vcsDirFor(repo *config.Repo) string {
	return fmt.Sprintf("vcs-%s", hashFor(repo.Url))
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Creates a new Searcher that is available for searches as soon as this returns.
// This will pull or clone the target repo and start watching the repo for changes.
func New(dbpath, name string, repo *config.Repo) (*Searcher, error) {
	vcsDir := filepath.Join(dbpath, vcsDirFor(repo))

	log.Printf("Searcher started for %s", name)

	sha, err := vcs.PullOrClone(repo.Vcs, vcsDir, repo.Url)
	if err != nil {
		return nil, err
	}

	idx, err := buildAndOpenIndex(dbpath, vcsDir)
	if err != nil {
		return nil, err
	}

	s := &Searcher{
		idx:  idx,
		Repo: repo,
	}

	go func() {
		for {
			time.Sleep(time.Duration(repo.MsBetweenPolls) * time.Millisecond)

			newSha, err := vcs.PullOrClone(repo.Vcs, vcsDir, repo.Url)
			if err != nil {
				log.Printf("vcs pull error (%s - %s): %s", name, repo.Url, err)
				continue
			}

			if newSha == sha {
				continue
			}

			log.Printf("Rebuilding %s for %s", name, newSha)
			idx, err := buildAndOpenIndex(dbpath, vcsDir)
			if err != nil {
				log.Printf("failed index build (%s): %s", name, err)
				os.RemoveAll(fmt.Sprintf("%s-%s", vcsDir, newSha))
				continue
			}

			if err := s.swapIndexes(idx); err != nil {
				log.Printf("failed index swap (%s): %s", name, err)
				if err := idx.Destroy(); err != nil {
					log.Printf("failed to destroy index (%s): %s\n", name, err)
				}
				continue
			}

			sha = newSha

			// This is just a good time to GC since we know there will be a
			// whole set of dead posting lists on the heap. Ensuring these
			// go away quickly helps to prevent the heap from expanding
			// uncessarily.
			runtime.GC()

			reportOnMemory()
		}
	}()

	return s, nil
}
