package client

import (
	"github.com/hound-search/hound/index"
)

type Block struct {
	Lines   []string
	Matches []bool
	Start   int
}

func endOfBlock(b *Block) int {
	return b.Start + len(b.Lines) - 1
}

func startOfMatch(m *index.Match) int {
	return m.LineNumber - len(m.Before)
}

func matchIsInBlock(m *index.Match, b *Block) bool {
	return startOfMatch(m) <= endOfBlock(b)
}

func matchToBlock(m *index.Match) *Block {
	b, a := len(m.Before), len(m.After)
	n := 1 + b + a
	l := make([]string, 0, n)
	v := make([]bool, n)

	v[b] = true

	for _, line := range m.Before {
		l = append(l, line)
	}

	l = append(l, m.Line)

	for _, line := range m.After {
		l = append(l, line)
	}

	return &Block{
		Lines:   l,
		Matches: v,
		Start:   m.LineNumber - len(m.Before),
	}
}

func clampZero(n int) int {
	if n < 0 {
		return 0
	}
	return n
}

func mergeMatchIntoBlock(m *index.Match, b *Block) {
	off := endOfBlock(b) - startOfMatch(m) + 1
	idx := len(b.Lines) - off
	nb := len(m.Before)

	for i := off; i < nb; i++ {
		b.Lines = append(b.Lines, m.Before[i])
		b.Matches = append(b.Matches, false)
	}

	if off < nb+1 {
		b.Lines = append(b.Lines, m.Line)
		b.Matches = append(b.Matches, true)
	} else {
		b.Matches[idx+nb] = true
	}

	for i, n := clampZero(off-nb-1), len(m.After); i < n; i++ {
		b.Lines = append(b.Lines, m.After[i])
		b.Matches = append(b.Matches, false)
	}
}

func coalesceMatches(matches []*index.Match) []*Block {
	var res []*Block
	var curr *Block
	for _, match := range matches {
		if curr != nil && matchIsInBlock(match, curr) {
			mergeMatchIntoBlock(match, curr)
		} else {
			if curr != nil {
				res = append(res, curr)
			}
			curr = matchToBlock(match)
		}
	}

	if curr != nil {
		res = append(res, curr)
	}

	return res
}
