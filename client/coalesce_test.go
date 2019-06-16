package client

import (
	"testing"

	"github.com/hound-search/hound/index"
)

// TODO(knorton):
// - Test multiple overlapping.
// - Test asymmetric context

func stringSlicesAreSame(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, n := 0, len(a); i < n; i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func boolSlicesAreSame(a, b []bool) bool {
	if len(a) != len(b) {
		return false
	}

	for i, n := 0, len(a); i < n; i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func assertBlocksAreSame(t *testing.T, a, b *Block) bool {
	if !stringSlicesAreSame(a.Lines, b.Lines) {
		t.Errorf("bad lines: expected: %v, got: %v", a.Lines, b.Lines)
		return false
	}

	if !boolSlicesAreSame(a.Matches, b.Matches) {
		t.Errorf("bad matches: expected: %v, got: %v", a.Matches, b.Matches)
		return false
	}

	if a.Start != b.Start {
		t.Errorf("bad start: expected %d, got %d", a.Start, b.Start)
		return false
	}

	return true
}

func assertBlockSlicesAreSame(t *testing.T, a, b []*Block) bool {
	if len(a) != len(b) {
		t.Errorf("blocks do not match, len(a)=%d & len(b)=%d", len(a), len(b))
		return false
	}

	for i, n := 0, len(a); i < n; i++ {
		if !assertBlocksAreSame(t, a[i], b[i]) {
			return false
		}
	}

	return true
}

func testThis(t *testing.T, subj []*index.Match, expt []*Block, desc string) {
	if !assertBlockSlicesAreSame(t, expt, coalesceMatches(subj)) {
		t.Errorf("case failed: %s", desc)
	}
}

func TestNonOverlap(t *testing.T) {
	subj := []*index.Match{
		&index.Match{
			Line:       "c",
			LineNumber: 40,
			Before:     []string{"a", "b"},
			After:      []string{"d", "e"},
		},
		&index.Match{
			Line:       "n",
			LineNumber: 50,
			Before:     []string{"l", "m"},
			After:      []string{"o", "p"},
		},
	}

	expt := []*Block{
		&Block{
			Lines:   []string{"a", "b", "c", "d", "e"},
			Matches: []bool{false, false, true, false, false},
			Start:   38,
		},
		&Block{
			Lines:   []string{"l", "m", "n", "o", "p"},
			Matches: []bool{false, false, true, false, false},
			Start:   48,
		},
	}

	testThis(t, subj, expt,
		"non-overlap w/ context")
}
func TestNonOverlapWithNoContext(t *testing.T) {
	subj := []*index.Match{
		&index.Match{
			Line:       "a",
			LineNumber: 40,
		},
		&index.Match{
			Line:       "b",
			LineNumber: 50,
		},
	}

	expt := []*Block{
		&Block{
			Lines:   []string{"a"},
			Matches: []bool{true},
			Start:   40,
		},

		&Block{
			Lines:   []string{"b"},
			Matches: []bool{true},
			Start:   50,
		},
	}

	testThis(t, subj, expt,
		"non-overlap w/o context")
}

func TestOverlappingInBefore(t *testing.T) {
	subj := []*index.Match{
		&index.Match{
			Line:       "c",
			LineNumber: 40,
			Before:     []string{"a", "b"},
			After:      []string{"d", "e"},
		},
		&index.Match{
			Line:       "g",
			LineNumber: 44,
			Before:     []string{"e", "f"},
			After:      []string{"h", "i"},
		},
	}

	expt := []*Block{
		&Block{
			Lines:   []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"},
			Matches: []bool{false, false, true, false, false, false, true, false, false},
			Start:   38,
		},
	}

	testThis(t, subj, expt,
		"overlap in before")
}
func TestOverlappingInAfter(t *testing.T) {
	subj := []*index.Match{
		&index.Match{
			Line:       "c",
			LineNumber: 40,
			Before:     []string{"a", "b"},
			After:      []string{"d", "e"},
		},
		&index.Match{
			Line:       "d",
			LineNumber: 41,
			Before:     []string{"b", "c"},
			After:      []string{"e", "f"},
		},
	}

	expt := []*Block{
		&Block{
			Lines:   []string{"a", "b", "c", "d", "e", "f"},
			Matches: []bool{false, false, true, true, false, false},
			Start:   38,
		},
	}

	testThis(t, subj, expt,
		"overlap in after")
}

func TestOverlapOnMatch(t *testing.T) {
	subj := []*index.Match{
		&index.Match{
			Line:       "c",
			LineNumber: 40,
			Before:     []string{"a", "b"},
			After:      []string{"d", "e"},
		},
		&index.Match{
			Line:       "e",
			LineNumber: 42,
			Before:     []string{"c", "d"},
			After:      []string{"f", "g"},
		},
	}

	expt := []*Block{
		&Block{
			Lines:   []string{"a", "b", "c", "d", "e", "f", "g"},
			Matches: []bool{false, false, true, false, true, false, false},
			Start:   38,
		},
	}

	testThis(t, subj, expt,
		"overlap on match")
}

func TestMatchesToEnd(t *testing.T) {
	file := []string{
		"import analytics.sequence._;",
		"import analytics._;",
		"println(\"Try running\")",
		"println(\"val visits = VisitExplorer(100)\");",
	}

	subj := []*index.Match{
		&index.Match{
			Line:       file[2],
			LineNumber: 3,
			Before:     []string{file[0], file[1]},
			After:      []string{file[3]},
		},

		&index.Match{
			Line:       file[3],
			LineNumber: 4,
			Before:     []string{file[1], file[2]},
			After:      nil,
		},
	}

	expt := []*Block{
		&Block{
			Lines:   []string{file[0], file[1], file[2], file[3]},
			Matches: []bool{false, false, true, true},
			Start:   1,
		},
	}

	testThis(t, subj, expt,
		"test matches at end of file")
}
