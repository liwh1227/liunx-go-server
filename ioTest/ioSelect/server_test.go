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
	type args struct {
		fd    int
		fdSet *syscall.FdSet
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FD_ISSET(tt.args.fd, tt.args.fdSet); got != tt.want {
				t.Errorf("FD_ISSET() = %v, want %v", got, tt.want)
			}
		})
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
