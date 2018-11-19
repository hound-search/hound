package ansi

// +build darwin freebsd openbsd netbsd dragonfly

import "syscall"

const ioctlReadTermios = syscall.TIOCGETA
