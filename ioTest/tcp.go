package ioTest

import (
	"fmt"
	"syscall"
)

// 建立tcp连接的bind、listen和socket描述符的一些设置
func PreStartTcp() (int, error) {
	//1.获取socketFd
	nSocketFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println("Get socket error.", err)
		return 0, err
	}
	//fmt.Printf("[nSocketFd] is %d\n",nSocketFd)
	//2. bind
	sa := &syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{0},
	}
	//2-1.设置socket描述符属性
	syscall.SetsockoptInt(nSocketFd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	err = syscall.Bind(nSocketFd, sa)
	if err != nil {
		fmt.Println("Bind error.", err)
		return 0, err
	}
	//2-2. listen socket
	err = syscall.Listen(nSocketFd, 1024)
	if err != nil {
		fmt.Printf("[listen] error :%s\n", err)
		return 0, err
	}
	return nSocketFd, nil
}
