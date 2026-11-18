//go:build !windows
// +build !windows

package ansi

import (
	"syscall"
	"unsafe"
)

// Issue a ioctl syscall to try to read a termios for the descriptor. If
// we are unable to read one, this is not a tty.
func isTTY(fd uintptr) bool {
	var termios syscall.Termios
	// OpenBSD routes SYS_IOCTL through libc only for the three-argument Syscall.
	_, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		ioctlReadTermios,
		uintptr(unsafe.Pointer(&termios)))
	return err == 0
}
