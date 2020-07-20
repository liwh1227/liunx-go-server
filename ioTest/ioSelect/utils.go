package ioSelect

import (
	"syscall"
)

func FD_ZERO(fdSet *syscall.FdSet) {
	for idx := range fdSet.Bits {
		fdSet.Bits[idx] = 0
	}
}

func FD_SET(fd int, fdSet *syscall.FdSet) {
	fdSet.Bits[(fd)>>5] |= 1 << ((fd) & 31)
}

func FD_CLR(fd int, fdSet *syscall.FdSet) {
	fdSet.Bits[(fd)>>5] &= ^(1 << ((fd) & 31))
}

func FD_ISSET(fd int, fdSet *syscall.FdSet) bool {
	return (fdSet.Bits[(fd)>>5] & (1 << ((fd) & 31))) != 0
}
