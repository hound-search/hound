package index

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/hound-search/hound/codesearch/regexp"
)

var (
	subjA = []byte("first\nsecond\nthird\nfourth\nfifth\nsixth")
	subjB = []byte("\n")
	subjC = []byte("\n\n\n\nfoo\nbar\n\nbaz")
)

func formatLines(lines []string) string {
	return fmt.Sprintf("[%s]", strings.Join(lines, ","))
}

func formatLinesFromBytes(lines [][]byte) string {
	n := len(lines)
	s := make([]string, n)
	for i := 0; i < n; i++ {
		s[i] = string(lines[i])
	}
	return formatLines(s)
}

func assertLinesMatch(t *testing.T, lines [][]byte, expected []string) {
	if len(lines) != len(expected) {
		t.Errorf("lines do not match: %s vs %s",
			formatLinesFromBytes(lines),
			formatLines(expected))
		return
	}
	for i, str := range expected {
		if str != string(lines[i]) {
			t.Errorf("lines do not match: %s vs %s",
				formatLinesFromBytes(lines),
				formatLines(expected))
		}
	}
}

func GenBuf(n int) []byte {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte(i)
	}
	return b
}

func TestFillFrom(t *testing.T) {
	var g grepper
	// this is to force buffer growth
	g.buf = make([]byte, 2)

	d := GenBuf(1024)
	b, _ := g.fillFrom(bytes.NewBuffer(d))
	if !bytes.Equal(d, b) {
		t.Errorf("filled buffer doesn't match original: len=%d & len=%d", len(d), len(b))
	}
}

func TestFirstNLines(t *testing.T) {
	assertLinesMatch(t, firstNLines(subjA, 1), []string{
		"first",
	})

	assertLinesMatch(t, firstNLines(subjA, 2), []string{
		"first",
		"second",
	})

	assertLinesMatch(t, firstNLines(subjA, 6), []string{
		"first",
		"second",
		"third",
		"fourth",
		"fifth",
		"sixth",
	})

	assertLinesMatch(t, firstNLines(subjB, 1), []string{
		"",
	})

	assertLinesMatch(t, firstNLines(subjB, 2), []string{
		"",
	})

	assertLinesMatch(t, firstNLines(subjC, 5), []string{
		"",
		"",
		"",
		"",
		"foo",
	})
}

func TestLastNLines(t *testing.T) {
	assertLinesMatch(t, lastNLines(subjA, 1), []string{
		"sixth",
	})

	assertLinesMatch(t, lastNLines(subjA, 2), []string{
		"fifth",
		"sixth",
	})

	assertLinesMatch(t, lastNLines(subjA, 6), []string{
		"first",
		"second",
		"third",
		"fourth",
		"fifth",
		"sixth",
	})

	assertLinesMatch(t, lastNLines(subjB, 1), []string{
		"",
	})

	assertLinesMatch(t, lastNLines(subjB, 2), []string{
		"",
	})

	assertLinesMatch(t, lastNLines(subjC, 5), []string{
		"",
		"foo",
		"bar",
		"",
		"baz",
	})
}

type match struct {
	line string
	no   int
}

func aMatch(line string, no int) *match {
	return &match{
		line: line,
		no:   no,
	}
}

func formatMatches(matches []*match) string {
	str := make([]string, len(matches))
	for i, match := range matches {
		str[i] = fmt.Sprintf("%s:%d", match.line, match.no)
	}
	return strings.Join(str, ",")
}

func assertMatchesMatch(t *testing.T, a, b []*match) {
	if len(a) != len(b) {
		t.Errorf("matches no match: %s vs %s",
			formatMatches(a),
			formatMatches(b))
		return
	}

	for i, n := 0, len(a); i < n; i++ {
		if a[i].line != b[i].line || a[i].no != b[i].no {
			t.Errorf("matches no match: %s vs %s",
				formatMatches(a),
				formatMatches(b))
			return
		}
	}
}

func assertGrepTest(t *testing.T, buf []byte, exp string, expects []*match) {
	re, err := regexp.Compile(exp)
	if err != nil {
		t.Error(err)
		return
	}

	var g grepper
	var m []*match
	if err := g.grep2(bytes.NewBuffer(buf), re, 0,
		func(line []byte, lineno int, before [][]byte, after [][]byte) (bool, error) {
			m = append(m, aMatch(string(line), lineno))
			return true, nil
		}); err != nil {
		t.Error(err)
		return
	}

	assertMatchesMatch(t, m, expects)
}

func TestGrep(t *testing.T) {
	assertGrepTest(t, subjA, "s", []*match{
		aMatch("first", 1),
		aMatch("second", 2),
		aMatch("sixth", 6),
	})

	// BUG(knorton): rsc's regexp has bugs.
	// assertGrepTest(t, subjB, "^$", []*match{
	//   aMatch("", 1),
	// })

	assertGrepTest(t, subjB, "^", []*match{
		aMatch("", 1),
	})

	assertGrepTest(t, subjC, "^", []*match{
		aMatch("", 1),
		aMatch("", 2),
		aMatch("", 3),
		aMatch("", 4),
		aMatch("foo", 5),
		aMatch("bar", 6),
		aMatch("", 7),
		aMatch("baz", 8),
	})

	assertGrepTest(t, subjA, "th$", []*match{
		aMatch("sixth", 6),
	})
}

func assertContextTest(t *testing.T, buf []byte, exp string, ctx int, expectsBefore [][]string, expectsAfter [][]string) {
	re, err := regexp.Compile(exp)
	if err != nil {
		t.Error(err)
		return
	}

	var gotBefore [][][]byte
	var gotAfter [][][]byte
	var g grepper
	if err := g.grep2(bytes.NewBuffer(buf), re, ctx,
		func(line []byte, lineno int, before [][]byte, after [][]byte) (bool, error) {
			gotBefore = append(gotBefore, before)
			gotAfter = append(gotAfter, after)
			return true, nil
		}); err != nil {
		t.Error(err)
		return
	}

	if len(expectsBefore) != len(gotBefore) {
		t.Errorf("Before had %d lines, should have had %d",
			len(gotBefore),
			len(expectsBefore))
		return
	}

	if len(expectsAfter) != len(gotAfter) {
		t.Errorf("After had %d lines, should have had %d",
			len(gotBefore),
			len(expectsBefore))
		return
	}

	for i, n := 0, len(gotBefore); i < n; i++ {
		assertLinesMatch(t, gotBefore[i], expectsBefore[i])
	}

	for i, n := 0, len(gotAfter); i < n; i++ {
		assertLinesMatch(t, gotAfter[i], expectsAfter[i])
	}
}

func TestContext(t *testing.T) {
	assertContextTest(t, subjA, "third", 2,
		[][]string{
			[]string{"first", "second"},
		}, [][]string{
			[]string{"fourth", "fifth"},
		})

	assertContextTest(t, subjA, "third", 3,
		[][]string{
			[]string{"first", "second"},
		}, [][]string{
			[]string{"fourth", "fifth", "sixth"},
		})

	assertContextTest(t, subjA, "first", 2,
		[][]string{
			[]string{},
		}, [][]string{
			[]string{"second", "third"},
		})
}
