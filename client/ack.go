package client

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/hound-search/hound/ansi"
	"github.com/hound-search/hound/config"
)

type ackPresenter struct {
	f *os.File
}

func hiliteMatches(c *ansi.Colorer, p *regexp.Regexp, line string) string {
	// find the indexes for all matches
	idxs := p.FindAllStringIndex(line, -1)

	var buf bytes.Buffer
	beg := 0

	for _, idx := range idxs {
		// for each match add the contents before the match ...
		buf.WriteString(line[beg:idx[0]])
		// and the highlighted version of the match
		buf.WriteString(c.FgBg(line[idx[0]:idx[1]],
			ansi.Black,
			ansi.Bold,
			ansi.Yellow,
			ansi.Intense))
		beg = idx[1]
	}

	buf.WriteString(line[beg:])

	return buf.String()
}

func lineNumber(c *ansi.Colorer, buf *bytes.Buffer, n int, hasMatch bool) string {
	defer buf.Reset()

	s := fmt.Sprintf("%d", n)
	buf.WriteString(c.Fg(s, ansi.Yellow, ansi.Bold))
	if hasMatch {
		buf.WriteByte(':')
	} else {
		buf.WriteByte('-')
	}
	for i := len(s); i < 6; i++ {
		buf.WriteByte(' ')
	}
	return buf.String()
}

func (p *ackPresenter) Present(
	re *regexp.Regexp,
	ctx int,
	repos map[string]*config.Repo,
	res *Response) error {

	c := ansi.NewFor(p.f)

	buf := bytes.NewBuffer(make([]byte, 0, 20))

	for repo, resp := range res.Results {
		if _, err := fmt.Fprintf(p.f, "%s\n",
			c.Fg(repoNameFor(repos, repo), ansi.Red, ansi.Bold)); err != nil {
			return err
		}

		for _, file := range resp.Matches {
			if _, err := fmt.Fprintf(p.f, "%s\n",
				c.Fg(file.Filename, ansi.Green, ansi.Bold)); err != nil {
				return err
			}

			blocks := coalesceMatches(file.Matches)

			for _, block := range blocks {
				for i, n := 0, len(block.Lines); i < n; i++ {
					line := block.Lines[i]
					hasMatch := block.Matches[i]

					if hasMatch {
						line = hiliteMatches(c, re, line)
					}

					if _, err := fmt.Fprintf(p.f, "%s%s\n",
						lineNumber(c, buf, block.Start+i, hasMatch),
						line); err != nil {
						return err
					}
				}

				if ctx > 0 {
					if _, err := fmt.Fprintln(p.f, "--"); err != nil {
						return err
					}
				}
			}

			if _, err := fmt.Fprintln(p.f); err != nil {
				return err
			}
		}
	}

	return nil
}

func NewAckPresenter(w *os.File) Presenter {
	return &ackPresenter{w}
}
