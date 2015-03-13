package index

import (
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/etsy/hound/codesearch/index"
	"github.com/etsy/hound/codesearch/regexp"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	matchLimit       = 5000
	manifestFilename = "metadata.gob"
)

type Index struct {
	Ref *IndexRef
	idx *index.Index
	lck sync.RWMutex
}

type SearchOptions struct {
	IgnoreCase     bool
	LinesOfContext uint
	FileRegexp     string
	Offset         int
	Limit          int
}

type Match struct {
	Line       string
	LineNumber int
	Before     []string
	After      []string
}

type SearchResponse struct {
	Matches        []*FileMatch
	FilesWithMatch int
	FilesOpened    int           `json:"-"`
	Duration       time.Duration `json:"-"`
	Revision       string
}

type FileMatch struct {
	Filename string
	Matches  []*Match
}

type ExcludedFile struct {
	Filename string
	Reason   string
}

type IndexRef struct {
	Url  string
	Rev  string
	Time time.Time
	dir  string
}

func (r *IndexRef) Dir() string {
	return r.dir
}

func (r *IndexRef) writeManifest() error {
	w, err := os.Create(filepath.Join(r.dir, manifestFilename))
	if err != nil {
		return err
	}
	defer w.Close()

	return gob.NewEncoder(w).Encode(r)
}

func (r *IndexRef) Open() (*Index, error) {
	return &Index{
		Ref: r,
		idx: index.Open(filepath.Join(r.dir, "tri")),
	}, nil
}

func (r *IndexRef) Remove() error {
	return os.RemoveAll(r.dir)
}

func (n *Index) Close() error {
	n.lck.Lock()
	defer n.lck.Unlock()
	return n.idx.Close()
}

func (n *Index) Destroy() error {
	n.lck.Lock()
	defer n.lck.Unlock()
	if err := n.idx.Close(); err != nil {
		return err
	}
	return n.Ref.Remove()
}

func (n *Index) GetDir() string {
	return n.Ref.dir
}

func toStrings(lines [][]byte) []string {
	strs := make([]string, len(lines))
	for i, n := 0, len(lines); i < n; i++ {
		strs[i] = string(lines[i])
	}
	return strs
}

func GetRegexpPattern(pat string, ignoreCase bool) string {
	if ignoreCase {
		return "(?i)(?m)" + pat
	}
	return "(?m)" + pat
}

func (n *Index) Search(pat string, opt *SearchOptions) (*SearchResponse, error) {
	startedAt := time.Now()

	n.lck.RLock()
	defer n.lck.RUnlock()

	re, err := regexp.Compile(GetRegexpPattern(pat, opt.IgnoreCase))
	if err != nil {
		return nil, err
	}

	var (
		g                grepper
		results          []*FileMatch
		filesOpened      int
		filesFound       int
		filesCollected   int
		matchesCollected int
	)

	var fre *regexp.Regexp
	if opt.FileRegexp != "" {
		fre, err = regexp.Compile(opt.FileRegexp)
		if err != nil {
			return nil, err
		}
	}

	files := n.idx.PostingQuery(index.RegexpQuery(re.Syntax))
	for _, file := range files {
		var matches []*Match
		name := n.idx.Name(file)
		hasMatch := false

		// reject files that do not match the file pattern
		if fre != nil && fre.MatchString(name, true, true) < 0 {
			continue
		}

		filesOpened++
		if err := g.grep2File(filepath.Join(n.Ref.dir, "raw", name), re, int(opt.LinesOfContext),
			func(line []byte, lineno int, before [][]byte, after [][]byte) (bool, error) {

				hasMatch = true
				if filesFound < opt.Offset || (opt.Limit > 0 && filesCollected >= opt.Limit) {
					return false, nil
				}

				matchesCollected++
				matches = append(matches, &Match{
					Line:       string(line),
					LineNumber: lineno,
					Before:     toStrings(before),
					After:      toStrings(after),
				})

				if matchesCollected > matchLimit {
					return false, fmt.Errorf("search exceeds limit on matches: %d", matchLimit)
				}

				return true, nil
			}); err != nil {
			return nil, err
		}

		if !hasMatch {
			continue
		}

		filesFound++
		if len(matches) > 0 {
			filesCollected++
			results = append(results, &FileMatch{
				Filename: name,
				Matches:  matches,
			})
		}
	}

	return &SearchResponse{
		Matches:        results,
		FilesWithMatch: filesFound,
		FilesOpened:    filesOpened,
		Duration:       time.Now().Sub(startedAt),
		Revision:       n.Ref.Rev,
	}, nil
}

func isTextFile(filename string) (bool, error) {
	buf := make([]byte, 2048)
	r, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer r.Close()

	n, err := io.ReadFull(r, buf)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return false, err
	}

	buf = buf[:n]

	return utf8.Valid(buf), nil
}

func addFileToIndex(ix *index.IndexWriter, dst, src, path string) (string, error) {
	rel, err := filepath.Rel(src, path)
	if err != nil {
		return "", err
	}

	r, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer r.Close()

	dup := filepath.Join(dst, "raw", rel)
	w, err := os.Create(dup)
	if err != nil {
		return "", err
	}
	defer w.Close()

	g := gzip.NewWriter(w)
	defer g.Close()

	return ix.Add(rel, io.TeeReader(r, g)), nil
}

func addDirToIndex(dst, src, path string) error {
	rel, err := filepath.Rel(src, path)
	if err != nil {
		return err
	}

	if rel == "." {
		return nil
	}

	dup := filepath.Join(dst, "raw", rel)
	return os.Mkdir(dup, os.ModePerm)
}

func indexAllFiles(dst, src string) error {
	ix := index.Create(filepath.Join(dst, "tri"))
	var excluded []*ExcludedFile

	// Make a file to store the excluded files for this repo
	fileHandle, err := os.Create(filepath.Join(dst, "excluded_files.json"))
	if err != nil {
		return err
	}
	defer fileHandle.Close()

	if err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		name := info.Name()
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		if name[0] == '.' {
			if info.IsDir() {
				return filepath.SkipDir
			}
			excluded = append(excluded, &ExcludedFile{rel, "Dot file"})
			return nil
		}

		if info.IsDir() {
			return addDirToIndex(dst, src, path)
		}

		if info.Mode()&os.ModeType != 0 {
			excluded = append(excluded, &ExcludedFile{rel, "Invalid Mode"})
			return nil
		}

		txt, err := isTextFile(path)
		if err != nil {
			return err
		}

		if !txt {
			excluded = append(excluded, &ExcludedFile{rel, "Not a text file"})
			return nil
		}

		reasonForExclusion, err := addFileToIndex(ix, dst, src, path)
		if err != nil {
			return err
		}
		if reasonForExclusion != "" {
			excluded = append(excluded, &ExcludedFile{rel, reasonForExclusion})
		}

		return nil
	}); err != nil {
		return err
	}

	if err := json.NewEncoder(fileHandle).Encode(excluded); err != nil {
		return err
	}

	ix.Flush()

	return nil
}

// Read the metadata for the index directory. Note that even if this
// returns a non-nil error, a Metadata object will be returned with
// all the information that is known about the index (this might
// include only the path)
func Read(dir string) (*IndexRef, error) {
	m := &IndexRef{
		dir: dir,
	}

	r, err := os.Open(filepath.Join(dir, manifestFilename))
	if err != nil {
		return m, err
	}
	defer r.Close()

	if err := gob.NewDecoder(r).Decode(m); err != nil {
		return m, err
	}

	return m, nil
}

func Build(dst, src, url, rev string) (*IndexRef, error) {
	if _, err := os.Stat(dst); err != nil {
		if err := os.MkdirAll(dst, os.ModePerm); err != nil {
			return nil, err
		}
	}

	if err := os.Mkdir(filepath.Join(dst, "raw"), os.ModePerm); err != nil {
		return nil, err
	}

	if err := indexAllFiles(dst, src); err != nil {
		return nil, err
	}

	r := &IndexRef{
		Url:  url,
		Rev:  rev,
		Time: time.Now(),
		dir:  dst,
	}

	if err := r.writeManifest(); err != nil {
		return nil, err
	}

	return r, nil
}

// Open the index in dir for searching.
func Open(dir string) (*Index, error) {
	r, err := Read(dir)
	if err != nil {
		return nil, err
	}

	return r.Open()
}
