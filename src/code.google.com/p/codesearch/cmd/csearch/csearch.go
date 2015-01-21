// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "runtime/pprof"

  "code.google.com/p/codesearch/index"
  "code.google.com/p/codesearch/regexp"
)

var usageMessage = `usage: csearch [-c] [-f fileregexp] [-h] [-i] [-l] [-n] regexp

Csearch behaves like grep over all indexed files, searching for regexp,
an RE2 (nearly PCRE) regular expression.

The -c, -h, -i, -l, and -n flags are as in grep, although note that as per Go's
flag parsing convention, they cannot be combined: the option pair -i -n 
cannot be abbreviated to -in.

The -f flag restricts the search to files whose names match the RE2 regular
expression fileregexp.

Csearch relies on the existence of an up-to-date index created ahead of time.
To build or rebuild the index that csearch uses, run:

	cindex path...

where path... is a list of directories or individual files to be included in the index.
If no index exists, this command creates one.  If an index already exists, cindex
overwrites it.  Run cindex -help for more.

Csearch uses the index stored in $CSEARCHINDEX or, if that variable is unset or
empty, $HOME/.csearchindex.
`

func usage() {
  fmt.Fprintf(os.Stderr, usageMessage)
  os.Exit(2)
}

var (
  fFlag       = flag.String("f", "", "search only files with names matching this regexp")
  iFlag       = flag.Bool("i", false, "case-insensitive search")
  verboseFlag = flag.Bool("verbose", false, "print extra information")
  bruteFlag   = flag.Bool("brute", false, "brute force - search all files in index")
  cpuProfile  = flag.String("cpuprofile", "", "write cpu profile to this file")

  matches bool
)

func Main() {
  g := regexp.Grep{
    Stdout: os.Stdout,
    Stderr: os.Stderr,
  }
  g.AddFlags()

  flag.Usage = usage
  flag.Parse()
  args := flag.Args()

  if len(args) != 1 {
    usage()
  }

  if *cpuProfile != "" {
    f, err := os.Create(*cpuProfile)
    if err != nil {
      log.Fatal(err)
    }
    defer f.Close()
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
  }

  pat := "(?m)" + args[0]
  if *iFlag {
    pat = "(?i)" + pat
  }
  re, err := regexp.Compile(pat)
  if err != nil {
    log.Fatal(err)
  }
  g.Regexp = re
  var fre *regexp.Regexp
  if *fFlag != "" {
    fre, err = regexp.Compile(*fFlag)
    if err != nil {
      log.Fatal(err)
    }
  }
  q := index.RegexpQuery(re.Syntax)
  if *verboseFlag {
    log.Printf("query: %s\n", q)
  }

  ix := index.Open(index.File())
  ix.Verbose = *verboseFlag
  var post []uint32
  if *bruteFlag {
    post = ix.PostingQuery(&index.Query{Op: index.QAll})
  } else {
    post = ix.PostingQuery(q)
  }
  if *verboseFlag {
    log.Printf("post query identified %d possible files\n", len(post))
  }

  if fre != nil {
    fnames := make([]uint32, 0, len(post))

    for _, fileid := range post {
      name := ix.Name(fileid)
      if fre.MatchString(name, true, true) < 0 {
        continue
      }
      fnames = append(fnames, fileid)
    }

    if *verboseFlag {
      log.Printf("filename regexp matched %d files\n", len(fnames))
    }
    post = fnames
  }

  for _, fileid := range post {
    name := ix.Name(fileid)
    g.File(name)
  }

  matches = g.Match
}

func main() {
  Main()
  if !matches {
    os.Exit(1)
  }
  os.Exit(0)
}
