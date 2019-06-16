package index

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"

	"github.com/hound-search/hound/codesearch/regexp"
)

var nl = []byte{'\n'}

type grepper struct {
	buf []byte
}

func countLines(b []byte) int {
	n := 0
	for {
		i := bytes.IndexByte(b, '\n')
		if i < 0 {
			break
		}
		n++
		b = b[i+1:]
	}
	return n
}

func (g *grepper) grepFile(filename string, re *regexp.Regexp,
	fn func(line []byte, lineno int) (bool, error)) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	c, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer c.Close()

	return g.grep(c, re, fn)
}

func (g *grepper) grep2File(filename string, re *regexp.Regexp, nctx int,
	fn func(line []byte, lineno int, before [][]byte, after [][]byte) (bool, error)) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	c, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer c.Close()

	return g.grep2(c, re, nctx, fn)
}

func (g *grepper) fillFrom(r io.Reader) ([]byte, error) {
	if g.buf == nil {
		g.buf = make([]byte, 1<<20)
	}

	off := 0
	for {
		n, err := io.ReadFull(r, g.buf[off:])
		if err == io.ErrUnexpectedEOF || err == io.EOF {
			return g.buf[:off+n], nil
		} else if err != nil {
			return nil, err
		}

		// grow the storage
		buf := make([]byte, len(g.buf)*2)
		copy(buf, g.buf)
		g.buf = buf
		off += n
	}
}

func lastNLines(buf []byte, n int) [][]byte {
	if len(buf) == 0 || n == 0 {
		return nil
	}

	r := make([][]byte, n)
	for i := 0; i < n; i++ {
		m := bytes.LastIndex(buf, nl)
		if m < 0 {
			if len(buf) == 0 {
				return r[n-i:]
			}
			r[n-i-1] = buf
			return r[n-i-1:]
		}
		r[n-i-1] = buf[m+1:]
		buf = buf[:m]
	}

	return r
}

func firstNLines(buf []byte, n int) [][]byte {
	if len(buf) == 0 || n == 0 {
		return nil
	}

	r := make([][]byte, n)
	for i := 0; i < n; i++ {
		m := bytes.Index(buf, nl)
		if m < 0 {
			if len(buf) == 0 {
				return r[:i]
			}
			r[i] = buf
			return r[:i+1]
		}
		r[i] = buf[:m]
		buf = buf[m+1:]
	}
	return r
}

// TODO(knorton): This is still being tested. This is a grep that supports context lines. Unlike the version
// in codesearch, this one does not operate on chunks. The downside is that we have to have the whole file
// in memory to do the grep. Fortunately, we limit the size of files that get indexed anyway. 10M files tend
// to not be source code.
func (g *grepper) grep2(
	r io.Reader,
	re *regexp.Regexp,
	nctx int,
	fn func(line []byte, lineno int, before [][]byte, after [][]byte) (bool, error)) error {

	buf, err := g.fillFrom(r)
	if err != nil {
		return err
	}

	lineno := 0
	for {
		if len(buf) == 0 {
			return nil
		}

		m := re.Match(buf, true, true)
		if m < 0 {
			return nil
		}

		// start of matched line.
		str := bytes.LastIndex(buf[:m], nl) + 1

		//end of previous line
		endl := str - 1
		if endl < 0 {
			endl = 0
		}

		//end of current line
		end := m + 1
		if end > len(buf) {
			end = len(buf)
		}

		lineno += countLines(buf[:str])

		more, err := fn(
			bytes.TrimRight(buf[str:end], "\n"),
			lineno+1,
			lastNLines(buf[:endl], nctx),
			firstNLines(buf[end:], nctx))
		if err != nil {
			return err
		}
		if !more {
			return nil
		}

		lineno++
		buf = buf[end:]
	}
}

// This nonsense is adapted from https://code.google.com/p/codesearch/source/browse/regexp/match.go#399
// and I assume it is a mess to make it faster, but I would like to try a much simpler cleaner version.
func (g *grepper) grep(r io.Reader, re *regexp.Regexp, fn func(line []byte, lineno int) (bool, error)) error {
	if g.buf == nil {
		g.buf = make([]byte, 1<<20)
	}

	var (
		buf       = g.buf[:0]
		lineno    = 1
		beginText = true
		endText   = false
	)

	for {
		n, err := io.ReadFull(r, buf[len(buf):cap(buf)])
		buf = buf[:len(buf)+n]
		end := len(buf)
		if err == nil {
			end = bytes.LastIndex(buf, nl) + 1
		} else {
			endText = true
		}
		chunkStart := 0
		for chunkStart < end {
			m1 := re.Match(buf[chunkStart:end], beginText, endText) + chunkStart
			beginText = false
			if m1 < chunkStart {
				break
			}
			lineStart := bytes.LastIndex(buf[chunkStart:m1], nl) + 1 + chunkStart
			lineEnd := m1 + 1
			if lineEnd > end {
				lineEnd = end
			}
			lineno += countLines(buf[chunkStart:lineStart])
			line := buf[lineStart:lineEnd]
			more, err := fn(line, lineno)
			if err != nil {
				return err
			}
			if !more {
				return nil
			}
			lineno++
			chunkStart = lineEnd
		}
		if err == nil {
			lineno += countLines(buf[chunkStart:end])
		}

		n = copy(buf, buf[end:])
		buf = buf[:n]
		if len(buf) == 0 && err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				return err
			}
			return nil
		}
	}

	return nil
}
