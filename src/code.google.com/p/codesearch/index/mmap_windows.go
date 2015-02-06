// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"log"
	"os"
	"syscall"
	"unsafe"
)

// An mmapData is mmap'ed read-only data from a file.
type mmapData struct {
  f *os.File
  d []byte
  o uintptr
}

func Munmap(o interface{}) error{
	// This Munmap function is called from read.go like Munmap(o unitptr) and
	// from write.go like Munmap(o []byte)
	// For now, only the calls from read.go are handled using the windows equivalent of Munmap
	// i.e syscall.UnmapViewOfFile(o uintptr)
	// Calls from write.go get a nil always.
	// TODO: Find a way to get the mmapData.o field in write.go.

	switch o.(type) {
		case uintptr:
			return syscall.UnmapViewOfFile(o.(uintptr))
		default:
			return nil	
	}
}

func mmapFile(f *os.File) mmapData {
	st, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	size := st.Size()
	if int64(int(size+4095)) != size+4095 {
		log.Fatalf("%s: too large for mmap", f.Name())
	}
	if size == 0 {
		var dummy uintptr
		return mmapData{f, nil, dummy}
	}
	h, err := syscall.CreateFileMapping(syscall.Handle(f.Fd()), nil, syscall.PAGE_READONLY, uint32(size>>32), uint32(size), nil)
	if err != nil {
		log.Fatalf("CreateFileMapping %s: %v", f.Name(), err)
	}

	addr, err := syscall.MapViewOfFile(h, syscall.FILE_MAP_READ, 0, 0, 0)
	if err != nil {
		log.Fatalf("MapViewOfFile %s: %v", f.Name(), err)
	}
	data := (*[1 << 30]byte)(unsafe.Pointer(addr))
	return mmapData{f, data[:size], addr}
}
