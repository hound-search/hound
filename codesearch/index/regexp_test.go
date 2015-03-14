// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"regexp/syntax"
	"testing"
)

var queryTests = []struct {
	re string
	q  string
}{
	{`Abcdef`, `"Abc" "bcd" "cde" "def"`},
	{`(abc)(def)`, `"abc" "bcd" "cde" "def"`},
	{`abc.*(def|ghi)`, `"abc" ("def"|"ghi")`},
	{`abc(def|ghi)`, `"abc" ("bcd" "cde" "def")|("bcg" "cgh" "ghi")`},
	{`a+hello`, `"ahe" "ell" "hel" "llo"`},
	{`(a+hello|b+world)`, `("ahe" "ell" "hel" "llo")|("bwo" "orl" "rld" "wor")`},
	{`a*bbb`, `"bbb"`},
	{`a?bbb`, `"bbb"`},
	{`(bbb)a?`, `"bbb"`},
	{`(bbb)a*`, `"bbb"`},
	{`^abc`, `"abc"`},
	{`abc$`, `"abc"`},
	{`ab[cde]f`, `("abc" "bcf")|("abd" "bdf")|("abe" "bef")`},
	{`(abc|bac)de`, `"cde" ("abc" "bcd")|("acd" "bac")`},

	// These don't have enough letters for a trigram, so they return the
	// always matching query "+".
	{`ab[^cde]f`, `+`},
	{`ab.f`, `+`},
	{`.`, `+`},
	{`()`, `+`},

	// No matches.
	{`[^\s\S]`, `-`},

	// Factoring works.
	{`(abc|abc)`, `"abc"`},
	{`(ab|ab)c`, `"abc"`},
	{`ab(cab|cat)`, `"abc" "bca" ("cab"|"cat")`},
	{`(z*(abc|def)z*)(z*(abc|def)z*)`, `("abc"|"def")`},
	{`(z*abcz*defz*)|(z*abcz*defz*)`, `"abc" "def"`},
	{`(z*abcz*defz*(ghi|jkl)z*)|(z*abcz*defz*(mno|prs)z*)`,
		`"abc" "def" ("ghi"|"jkl"|"mno"|"prs")`},
	{`(z*(abcz*def)|(ghiz*jkl)z*)|(z*(mnoz*prs)|(tuvz*wxy)z*)`,
		`("abc" "def")|("ghi" "jkl")|("mno" "prs")|("tuv" "wxy")`},
	{`(z*abcz*defz*)(z*(ghi|jkl)z*)`, `"abc" "def" ("ghi"|"jkl")`},
	{`(z*abcz*defz*)|(z*(ghi|jkl)z*)`, `("ghi"|"jkl")|("abc" "def")`},

	// analyze keeps track of multiple possible prefix/suffixes.
	{`[ab][cd][ef]`, `("ace"|"acf"|"ade"|"adf"|"bce"|"bcf"|"bde"|"bdf")`},
	{`ab[cd]e`, `("abc" "bce")|("abd" "bde")`},

	// Different sized suffixes.
	{`(a|ab)cde`, `"cde" ("abc" "bcd")|("acd")`},
	{`(a|b|c|d)(ef|g|hi|j)`, `+`},

	{`(?s).`, `+`},

	// Expanding case.
	{`(?i)a~~`, `("A~~"|"a~~")`},
	{`(?i)ab~`, `("AB~"|"Ab~"|"aB~"|"ab~")`},
	{`(?i)abc`, `("ABC"|"ABc"|"AbC"|"Abc"|"aBC"|"aBc"|"abC"|"abc")`},
	{`(?i)abc|def`, `("ABC"|"ABc"|"AbC"|"Abc"|"DEF"|"DEf"|"DeF"|"Def"|"aBC"|"aBc"|"abC"|"abc"|"dEF"|"dEf"|"deF"|"def")`},
	{`(?i)abcd`, `("ABC"|"ABc"|"AbC"|"Abc"|"aBC"|"aBc"|"abC"|"abc") ("BCD"|"BCd"|"BcD"|"Bcd"|"bCD"|"bCd"|"bcD"|"bcd")`},
	{`(?i)abc|abc`, `("ABC"|"ABc"|"AbC"|"Abc"|"aBC"|"aBc"|"abC"|"abc")`},

	// Word boundary.
	{`\b`, `+`},
	{`\B`, `+`},
	{`\babc`, `"abc"`},
	{`\Babc`, `"abc"`},
	{`abc\b`, `"abc"`},
	{`abc\B`, `"abc"`},
	{`ab\bc`, `"abc"`},
	{`ab\Bc`, `"abc"`},
}

func TestQuery(t *testing.T) {
	for _, tt := range queryTests {
		re, err := syntax.Parse(tt.re, syntax.Perl)
		if err != nil {
			t.Fatal(err)
		}
		q := RegexpQuery(re).String()
		if q != tt.q {
			t.Errorf("RegexpQuery(%#q) = %#q, want %#q", tt.re, q, tt.q)
		}
	}
}
