package ioSelect

import (
	"syscall"
	"testing"
)

func TestFD_CLR(t *testing.T) {
	fdSet := syscall.FdSet{Bits: [16]int64{1}}
	FD_SET(2, (*syscall.FdSet)(&fdSet))
	t.Log(fdSet)
}

func TestFD_ISSET(t *testing.T) {
	fdSet := syscall.FdSet{Bits: [16]int64{1}}
	got := FD_ISSET(5, &fdSet)
	if got {
		t.Log("is True")
	} else {
		t.Error("fd_isset error")
	}

}

func TestFD_SET(t *testing.T) {
	type args struct {
		fd    int
		fdSet *syscall.FdSet
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestFD_ZERO(t *testing.T) {
	type args struct {
		fdSet *syscall.FdSet
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
