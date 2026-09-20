// +build darwin freebsd openbsd netbsd dragonfly
package ansi

import "syscall"

const ioctlReadTermios = syscall.TIOCGETA
