// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package regexp

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

var nstateTests = []struct {
	q       []uint32
	partial rune
}{
	{[]uint32{1, 2, 3}, 1},
	{[]uint32{1}, 1},
	{[]uint32{}, 0},
	{[]uint32{1, 2, 8}, 0x10FFF},
}

func TestNstateEnc(t *testing.T) {
	var n1, n2 nstate
	n1.q.Init(10)
	n2.q.Init(10)
	for _, tt := range nstateTests {
		n1.q.Reset()
		n1.partial = tt.partial
		for _, id := range tt.q {
			n1.q.Add(id)
		}
		enc := n1.enc()
		n2.dec(enc)
		if n2.partial != n1.partial || !reflect.DeepEqual(n1.q.Dense(), n2.q.Dense()) {
			t.Errorf("%v.enc.dec = %v", &n1, &n2)
		}
	}
}

var matchTests = []struct {
	re string
	s  string
	m  []int
}{
	// Adapted from go/src/pkg/regexp/find_test.go.
	{`a+`, "abc\ndef\nghi\n", []int{1}},
	{``, ``, []int{1}},
	{`^abcdefg`, "abcdefg", []int{1}},
	{`a+`, "baaab", []int{1}},
	{"abcd..", "abcdef", []int{1}},
	{`a`, "a", []int{1}},
	{`x`, "y", nil},
	{`b`, "abc", []int{1}},
	{`.`, "a", []int{1}},
	{`.*`, "abcdef", []int{1}},
	{`^`, "abcde", []int{1}},
	{`$`, "abcde", []int{1}},
	{`^abcd$`, "abcd", []int{1}},
	{`^bcd'`, "abcdef", nil},
	{`^abcd$`, "abcde", nil},
	{`a+`, "baaab", []int{1}},
	{`a*`, "baaab", []int{1}},
	{`[a-z]+`, "abcd", []int{1}},
	{`[^a-z]+`, "ab1234cd", []int{1}},
	{`[a\-\]z]+`, "az]-bcz", []int{1}},
	{`[^\n]+`, "abcd\n", []int{1}},
	{`[日本語]+`, "日本語日本語", []int{1}},
	{`日本語+`, "日本語", []int{1}},
	{`日本語+`, "日本語語語語", []int{1}},
	{`()`, "", []int{1}},
	{`(a)`, "a", []int{1}},
	{`(.)(.)`, "日a", []int{1}},
	{`(.*)`, "", []int{1}},
	{`(.*)`, "abcd", []int{1}},
	{`(..)(..)`, "abcd", []int{1}},
	{`(([^xyz]*)(d))`, "abcd", []int{1}},
	{`((a|b|c)*(d))`, "abcd", []int{1}},
	{`(((a|b|c)*)(d))`, "abcd", []int{1}},
	{`\a\f\r\t\v`, "\a\f\r\t\v", []int{1}},
	{`[\a\f\n\r\t\v]+`, "\a\f\r\t\v", []int{1}},

	{`a*(|(b))c*`, "aacc", []int{1}},
	{`(.*).*`, "ab", []int{1}},
	{`[.]`, ".", []int{1}},
	{`/$`, "/abc/", []int{1}},
	{`/$`, "/abc", nil},

	// multiple matches
	{`.`, "abc", []int{1}},
	{`(.)`, "abc", []int{1}},
	{`.(.)`, "abcd", []int{1}},
	{`ab*`, "abbaab", []int{1}},
	{`a(b*)`, "abbaab", []int{1}},

	// fixed bugs
	{`ab$`, "cab", []int{1}},
	{`axxb$`, "axxcb", nil},
	{`data`, "daXY data", []int{1}},
	{`da(.)a$`, "daXY data", []int{1}},
	{`zx+`, "zzx", []int{1}},
	{`ab$`, "abcab", []int{1}},
	{`(aa)*$`, "a", []int{1}},
	{`(?:.|(?:.a))`, "", nil},
	{`(?:A(?:A|a))`, "Aa", []int{1}},
	{`(?:A|(?:A|a))`, "a", []int{1}},
	{`(a){0}`, "", []int{1}},
	//	{`(?-s)(?:(?:^).)`, "\n", nil},
	//	{`(?s)(?:(?:^).)`, "\n", []int{1}},
	//	{`(?:(?:^).)`, "\n", nil},
	{`\b`, "x", []int{1}},
	{`\b`, "xx", []int{1}},
	{`\b`, "x y", []int{1}},
	{`\b`, "xx yy", []int{1}},
	{`\B`, "x", nil},
	{`\B`, "xx", []int{1}},
	{`\B`, "x y", nil},
	{`\B`, "xx yy", []int{1}},
	{`(?im)^[abc]+$`, "abcABC", []int{1}},
	{`(?im)^[α]+$`, "αΑ", []int{1}},
	{`[Aa]BC`, "abc", nil},
	{`[Aa]bc`, "abc", []int{1}},

	// RE2 tests
	{`[^\S\s]`, "abcd", nil},
	{`[^\S[:space:]]`, "abcd", nil},
	{`[^\D\d]`, "abcd", nil},
	{`[^\D[:digit:]]`, "abcd", nil},
	{`(?i)\W`, "x", nil},
	{`(?i)\W`, "k", nil},
	{`(?i)\W`, "s", nil},

	// can backslash-escape any punctuation
	{`\!\"\#\$\%\&\'\(\)\*\+\,\-\.\/\:\;\<\=\>\?\@\[\\\]\^\_\{\|\}\~`,
		`!"#$%&'()*+,-./:;<=>?@[\]^_{|}~`, []int{1}},
	{`[\!\"\#\$\%\&\'\(\)\*\+\,\-\.\/\:\;\<\=\>\?\@\[\\\]\^\_\{\|\}\~]+`,
		`!"#$%&'()*+,-./:;<=>?@[\]^_{|}~`, []int{1}},
	{"\\`", "`", []int{1}},
	{"[\\`]+", "`", []int{1}},

	// long set of matches (longer than startSize)
	{
		".",
		"qwertyuiopasdfghjklzxcvbnm1234567890",
		[]int{1},
	},
}

func TestMatch(t *testing.T) {
	for _, tt := range matchTests {
		re, err := Compile("(?m)" + tt.re)
		if err != nil {
			t.Errorf("Compile(%#q): %v", tt.re, err)
			continue
		}
		b := []byte(tt.s)
		lines := grep(re, b)
		if !reflect.DeepEqual(lines, tt.m) {
			t.Errorf("grep(%#q, %q) = %v, want %v", tt.re, tt.s, lines, tt.m)
		}
	}
}

func grep(re *Regexp, b []byte) []int {
	var m []int
	lineno := 1
	for {
		i := re.Match(b, true, true)
		if i < 0 {
			break
		}
		start := bytes.LastIndex(b[:i], nl) + 1
		end := i + 1
		if end > len(b) {
			end = len(b)
		}
		lineno += bytes.Count(b[:start], nl)
		m = append(m, lineno)
		if start < end && b[end-1] == '\n' {
			lineno++
		}
		b = b[end:]
		if len(b) == 0 {
			break
		}
	}
	return m
}

var grepTests = []struct {
	re  string
	s   string
	out string
	err string
	g   Grep
}{
	{re: `a+`, s: "abc\ndef\nghalloo\n", out: "input:abc\ninput:ghalloo\n"},
	{re: `x.*y`, s: "xay\nxa\ny\n", out: "input:xay\n"},
}

func TestGrep(t *testing.T) {
	for i, tt := range grepTests {
		re, err := Compile("(?m)" + tt.re)
		if err != nil {
			t.Errorf("Compile(%#q): %v", tt.re, err)
			continue
		}
		g := tt.g
		g.Regexp = re
		var out, errb bytes.Buffer
		g.Stdout = &out
		g.Stderr = &errb
		g.Reader(strings.NewReader(tt.s), "input")
		if out.String() != tt.out || errb.String() != tt.err {
			t.Errorf("#%d: grep(%#q, %q) = %q, %q, want %q, %q", i, tt.re, tt.s, out.String(), errb.String(), tt.out, tt.err)
		}
	}
}
