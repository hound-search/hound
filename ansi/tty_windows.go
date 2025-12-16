package ansi

import (
	"syscall"
	"unsafe"
)

var (
	modkernel32        = syscall.MustLoadDLL("kernel32.dll")
	procGetConsoleMode = modkernel32.MustFindProc("GetConsoleMode")
)

func isTTY(fd uintptr) bool {
	var mode uint32
	ret, _, err := procGetConsoleMode.Call(fd, uintptr(unsafe.Pointer(&mode)))
	return ret != 0 && err != nil
}
