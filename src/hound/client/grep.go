package client

import (
	"fmt"
	"hound/config"
	"os"
	"regexp"
)

type grepPresenter struct {
	f *os.File
}

func (p *grepPresenter) Present(
	re *regexp.Regexp,
	ctx int,
	repos map[string]*config.Repo,
	res *Response) error {

	for repo, resp := range res.Results {
		for _, file := range resp.Matches {
			for _, match := range file.Matches {
				if _, err := fmt.Fprintf(p.f, "%s/%s:%d: %s\n",
					repoNameFor(repos, repo),
					file.Filename,
					match.LineNumber,
					match.Line); err != nil {
						return err
					}
			}
		}
	}

	return nil
}

func NewGrepPresenter(w *os.File) Presenter {
	return &grepPresenter{w}
}
