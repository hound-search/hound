package ansi

import (
	"syscall"
	"unsafe"
)

// Issue a ioctl syscall to try to read a termios for the descriptor. If
// we are unable to read one, this is not a tty.
func isTTY(fd uintptr) bool {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(
		syscall.SYS_IOCTL,
		fd,
		ioctlReadTermios,
		uintptr(unsafe.Pointer(&termios)),
		0,
		0,
		0)
	return err == 0
}
