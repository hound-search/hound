package searcher

import (
	"fmt"
	"hound/config"
	"hound/index"
	"hound/vcs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
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

func expungeOldIndexes(sha, vcsDir string) error {
	// TODO(knorton): This is a bandaid for issue #14, but should suffice
	// since people don't usually name their repos with 40 char hashes. In
	// the longer term, I want to remove this naming scheme to support
	// rebuilds of the current hash.
	pat := regexp.MustCompile("-[0-9a-f]{40}$")

	name := fmt.Sprintf("%s-%s", filepath.Base(vcsDir), sha)

	dirs, err := filepath.Glob(fmt.Sprintf("%s-*", vcsDir))
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		bn := filepath.Base(dir)
		if !pat.MatchString(bn) || len(bn) != len(name) {
			continue
		}

		if bn == name {
			continue
		}

		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}

	return nil
}

func buildAndOpenIndex(sha, vcsDir string) (*index.Index, error) {
	idxDir := fmt.Sprintf("%s-%s", vcsDir, sha)
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

// Creates a new Searcher for the vcsDir but avoids any remote vcs operations.
// This requires that an existing vcsDir be available in the data directory. This
// is intended for debugging and testing only. This will not start a watcher to
// monitor the remote repo for changes.
func NewFromExisting(vcsDir string, repo *config.Repo) (*Searcher, error) {
	name := filepath.Base(vcsDir)

	log.Printf("Search started for %s", name)
	log.Println("  WARNING: index is static and will not update")

	sha, err := vcs.HeadHash(repo.VCS, vcsDir)
	if err != nil {
		return nil, err
	}

	idx, err := buildAndOpenIndex(sha, vcsDir)
	if err != nil {
		return nil, err
	}

	return &Searcher{
		idx:  idx,
		Repo: repo,
	}, nil
}

// Creates a new Searcher that is available for searches as soon as this returns.
// This will pull or clone the target repo and start watching the repo for changes.
func New(vcsDir string, repo *config.Repo) (*Searcher, error) {
	name := filepath.Base(vcsDir)

	log.Printf("Searcher started for %s", name)

	sha, err := vcs.PullOrClone(repo.VCS, vcsDir, repo.Url)
	if err != nil {
		return nil, err
	}

	if err := expungeOldIndexes(sha, vcsDir); err != nil {
		return nil, err
	}

	idx, err := buildAndOpenIndex(sha, vcsDir)
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

			newSha, err := vcs.PullOrClone(repo.VCS, vcsDir, repo.Url)
			if err != nil {
				log.Printf("vcs pull error (%s - %s): %s", name, repo.Url, err)
				continue
			}

			if newSha == sha {
				continue
			}

			log.Printf("Rebuilding %s for %s", name, newSha)
			idx, err := buildAndOpenIndex(newSha, vcsDir)
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
