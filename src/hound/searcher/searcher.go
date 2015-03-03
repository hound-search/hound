package searcher

import (
	"crypto/sha1"
	"encoding/gob"
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
	"strings"
	"sync"
	"time"
)

const (
	manifestFileName = "manifest.gob"
)

type Searcher struct {
	idx  *index.Index
	lck  sync.RWMutex
	Repo *config.Repo
}

type manifest struct {
	Url  string
	Rev  string
	Time time.Time

	path string
	keep bool
}

// Perform atomic swap of index in the searcher so that the new
// index is made "live".
func (s *Searcher) swapIndexes(idx *index.Index) error {
	s.lck.Lock()
	defer s.lck.Unlock()

	oldIdx := s.idx
	s.idx = idx

	return oldIdx.Destroy()
}

// Perform a basic search on the current index using the supplied pattern
// and the options.
//
// TODO(knorton): pat should really just be a part of SearchOptions
func (s *Searcher) Search(pat string, opt *index.SearchOptions) (*index.SearchResponse, error) {
	s.lck.RLock()
	defer s.lck.RUnlock()
	return s.idx.Search(pat, opt)
}

// Get the excluded files as a JSON string. This is only used for returning
// the data directly to clients (thus JSON).
func (s *Searcher) GetExcludedFiles() string {
	path := filepath.Join(s.idx.GetDir(), "excluded_files.json")
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Couldn't read excluded_files.json %v\n", err)
	}
	return string(dat)
}

// Generate a new index directory in the dbpath. The names are based
// on pseudo-randomness with a time-based seed.
func nextIndexDir(dbpath string) string {
	r := uint64(rand.Uint32())<<32 | uint64(rand.Uint32())
	return filepath.Join(dbpath, fmt.Sprintf("idx-%08x", r))
}

// Read the manifest for the index directory. Note that even if
// this returns a non-nil error, a manifest will be returned with
// all the information we do know about the index dir. The most
// common error condition will return a manifest with only a path
// set. This allows us to clean up dead, mal-formed indexes.
func readManifest(idxDir string) (*manifest, error) {
	m := &manifest{
		path: idxDir,
	}

	r, err := os.Open(filepath.Join(idxDir, manifestFileName))
	if err != nil {
		return m, err
	}
	defer r.Close()

	if err := gob.NewDecoder(r).Decode(m); err != nil {
		return m, err
	}

	return m, nil
}

// Serialize this manifest into the given index directory.
func (m *manifest) write(idxDir string) error {
	w, err := os.Create(filepath.Join(idxDir, manifestFileName))
	if err != nil {
		return err
	}
	defer w.Close()

	return gob.NewEncoder(w).Encode(m)
}

// Remove the index directory associated with this manifest.
func (m *manifest) removeIndex() error {
	return os.RemoveAll(m.path)
}

// Read the manifests associated with each of the index dirs
// in the given dbpath.
func readAllManifests(dbpath string) ([]*manifest, error) {
	dirs, err := filepath.Glob(filepath.Join(dbpath, "idx-*"))
	if err != nil {
		return nil, err
	}

	var ms []*manifest
	for _, dir := range dirs {
		m, _ := readManifest(dir)
		ms = append(ms, m)
	}

	return ms, nil
}

// Find a manifest corresponding to a particular repo and revision. This
// returns nil if no such manifest exists. Note that this is simple linear
// search but we are typically dealing with small numbers of indices, so the
// cost will be less than hashing string keys.
func findManifest(manifests []*manifest, url, rev string) *manifest {
	for _, m := range manifests {
		if m.Url == url && m.Rev == rev {
			return m
		}
	}
	return nil
}

// Remove all indexes associated with manfiests that have not been
// kept.
func removeUnusedIndexes(manifests []*manifest) error {
	for _, m := range manifests {
		if m.keep {
			continue
		}

		if err := m.removeIndex(); err != nil {
			return err
		}
	}

	return nil
}

// Open an index at the given path. If the idxDir is already present, it will
// simply open and use that index. If, however, the idxDir does not exist a new
// one will be built.
func buildAndOpenIndex(dbpath, vcsDir, idxDir, url, rev string) (*index.Index, error) {
	if _, err := os.Stat(idxDir); err != nil {
		_, err := index.Build(idxDir, vcsDir)
		if err != nil {
			return nil, err
		}

		m := manifest{
			Url:  url,
			Rev:  rev,
			Time: time.Now(),
		}

		if err := m.write(idxDir); err != nil {
			return nil, err
		}
	}

	return index.Open(idxDir)
}

// Simply prints out statistics about the heap. When hound rebuilds a new
// index it will expand the heap with a decent amount of garbage. This is
// helpful to ensure the heap growth looks sane.
func reportOnMemory() {
	var ms runtime.MemStats

	// Print out interesting heap info.
	runtime.ReadMemStats(&ms)
	fmt.Printf("HeapInUse = %0.2f\n", float64(ms.HeapInuse)/1e6)
	fmt.Printf("HeapIdle  = %0.2f\n", float64(ms.HeapIdle)/1e6)
}

// Utility function for producing a hex encoded sha1 hash for a string.
func hashFor(name string) string {
	h := sha1.New()
	h.Write([]byte(name))
	return hex.EncodeToString(h.Sum(nil))
}

// Create a normalized name for the vcs directory of this repo.
func vcsDirFor(repo *config.Repo) string {
	return fmt.Sprintf("vcs-%s", hashFor(repo.Url))
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Make a searcher for each repo in the Config. This function kind of has a notion
// of partial errors. First, if the error returned is non-nil then a fatal error has
// occurred and no other return values are valid. If an error occurs that is specific
// to a particular searcher, that searcher will not be present in the searcher map and
// will have an error entry in the error map.
//
// NOTE: The keys in the searcher map will be normalized to lower case, but not such
// transformation will be done on the error map to make it easier to match those errors
// back to the original repo name.
func MakeAll(cfg *config.Config) (map[string]*Searcher, map[string]error, error) {
	errs := map[string]error{}
	searchers := map[string]*Searcher{}

	manifests, err := readAllManifests(cfg.DbPath)
	if err != nil {
		return nil, nil, err
	}

	for name, repo := range cfg.Repos {
		s, err := newSearcher(cfg.DbPath, name, repo, manifests)
		if err != nil {
			errs[name] = err
		}

		searchers[strings.ToLower(name)] = s
	}

	if err := removeUnusedIndexes(manifests); err != nil {
		return nil, nil, err
	}

	return searchers, errs, nil
}

// Creates a new Searcher that is available for searches as soon as this returns.
// This will pull or clone the target repo and start watching the repo for changes.
func New(dbpath, name string, repo *config.Repo) (*Searcher, error) {
	return newSearcher(dbpath, name, repo, nil)
}

// Creates a new Searcher that is capable of re-claiming an existing index directory
// from a set of existing manifests.
func newSearcher(dbpath, name string, repo *config.Repo, manifests []*manifest) (*Searcher, error) {
	vcsDir := filepath.Join(dbpath, vcsDirFor(repo))

	log.Printf("Searcher started for %s", name)

	rev, err := vcs.PullOrClone(repo.Vcs, vcsDir, repo.Url)
	if err != nil {
		return nil, err
	}

	var idxDir string
	man := findManifest(manifests, repo.Url, rev)
	if man == nil {
		idxDir = nextIndexDir(dbpath)
	} else {
		idxDir = man.path
		man.keep = true
	}

	idx, err := buildAndOpenIndex(dbpath, vcsDir, idxDir, repo.Url, rev)
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

			newRev, err := vcs.PullOrClone(repo.Vcs, vcsDir, repo.Url)
			if err != nil {
				log.Printf("vcs pull error (%s - %s): %s", name, repo.Url, err)
				continue
			}

			if newRev == rev {
				continue
			}

			log.Printf("Rebuilding %s for %s", name, newRev)
			idx, err := buildAndOpenIndex(
				dbpath,
				vcsDir,
				nextIndexDir(dbpath),
				repo.Url,
				newRev)
			if err != nil {
				log.Printf("failed index build (%s): %s", name, err)
				continue
			}

			if err := s.swapIndexes(idx); err != nil {
				log.Printf("failed index swap (%s): %s", name, err)
				if err := idx.Destroy(); err != nil {
					log.Printf("failed to destroy index (%s): %s\n", name, err)
				}
				continue
			}

			rev = newRev

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
