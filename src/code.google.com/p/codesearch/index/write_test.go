// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"bytes"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"testing"
)

var trivialFiles = map[string]string{
	"f0":       "\n\n",
	"file1":    "\na\n",
	"thefile2": "\nab\n",
	"file3":    "\nabc\n",
	"afile4":   "\ndabc\n",
	"file5":    "\nxyzw\n",
}

var trivialIndex = join(
	// header
	"csearch index 1\n",

	// list of paths
	"\x00",

	// list of names
	"afile4\x00",
	"f0\x00",
	"file1\x00",
	"file3\x00",
	"file5\x00",
	"thefile2\x00",
	"\x00",

	// list of posting lists
	"\na\n", fileList(2), // file1
	"\nab", fileList(3, 5), // file3, thefile2
	"\nda", fileList(0), // afile4
	"\nxy", fileList(4), // file5
	"ab\n", fileList(5), // thefile2
	"abc", fileList(0, 3), // afile4, file3
	"bc\n", fileList(0, 3), // afile4, file3
	"dab", fileList(0), // afile4
	"xyz", fileList(4), // file5
	"yzw", fileList(4), // file5
	"zw\n", fileList(4), // file5
	"\xff\xff\xff", fileList(),

	// name index
	u32(0),
	u32(6+1),
	u32(6+1+2+1),
	u32(6+1+2+1+5+1),
	u32(6+1+2+1+5+1+5+1),
	u32(6+1+2+1+5+1+5+1+5+1),
	u32(6+1+2+1+5+1+5+1+5+1+8+1),

	// posting list index,
	"\na\n", u32(1), u32(0),
	"\nab", u32(2), u32(5),
	"\nda", u32(1), u32(5+6),
	"\nxy", u32(1), u32(5+6+5),
	"ab\n", u32(1), u32(5+6+5+5),
	"abc", u32(2), u32(5+6+5+5+5),
	"bc\n", u32(2), u32(5+6+5+5+5+6),
	"dab", u32(1), u32(5+6+5+5+5+6+6),
	"xyz", u32(1), u32(5+6+5+5+5+6+6+5),
	"yzw", u32(1), u32(5+6+5+5+5+6+6+5+5),
	"zw\n", u32(1), u32(5+6+5+5+5+6+6+5+5+5),
	"\xff\xff\xff", u32(0), u32(5+6+5+5+5+6+6+5+5+5+5),

	// trailer
	u32(16),
	u32(16+1),
	u32(16+1+38),
	u32(16+1+38+62),
	u32(16+1+38+62+28),

	"\ncsearch trailr\n",
)

func join(s ...string) string {
	return strings.Join(s, "")
}

func u32(x uint32) string {
	var buf [4]byte
	buf[0] = byte(x >> 24)
	buf[1] = byte(x >> 16)
	buf[2] = byte(x >> 8)
	buf[3] = byte(x)
	return string(buf[:])
}

func fileList(list ...uint32) string {
	var buf []byte

	last := ^uint32(0)
	for _, x := range list {
		delta := x - last
		for delta >= 0x80 {
			buf = append(buf, byte(delta)|0x80)
			delta >>= 7
		}
		buf = append(buf, byte(delta))
		last = x
	}
	buf = append(buf, 0)
	return string(buf)
}

func buildFlushIndex(out string, paths []string, doFlush bool, fileData map[string]string) {
	ix := Create(out)
	ix.AddPaths(paths)
	var files []string
	for name := range fileData {
		files = append(files, name)
	}
	sort.Strings(files)
	for _, name := range files {
		ix.Add(name, strings.NewReader(fileData[name]))
	}
	if doFlush {
		ix.flushPost()
	}
	ix.Flush()
}

func buildIndex(name string, paths []string, fileData map[string]string) {
	buildFlushIndex(name, paths, false, fileData)
}

func testTrivialWrite(t *testing.T, doFlush bool) {
	f, _ := ioutil.TempFile("", "index-test")
	defer os.Remove(f.Name())
	out := f.Name()
	buildFlushIndex(out, nil, doFlush, trivialFiles)

	data, err := ioutil.ReadFile(out)
	if err != nil {
		t.Fatalf("reading _test/index.triv: %v", err)
	}
	want := []byte(trivialIndex)
	if !bytes.Equal(data, want) {
		i := 0
		for i < len(data) && i < len(want) && data[i] == want[i] {
			i++
		}
		t.Fatalf("wrong index:\nhave: %q %q\nwant: %q %q", data[:i], data[i:], want[:i], want[i:])
	}
}

func TestTrivialWrite(t *testing.T) {
	testTrivialWrite(t, false)
}

func TestTrivialWriteDisk(t *testing.T) {
	testTrivialWrite(t, true)
}
