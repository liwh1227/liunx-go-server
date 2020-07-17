package ioSelect

/*
   模拟linux select中的方法:
   void FD_ZERO (fd_set);
   //设置fdset中需要监控的bit位
   void FD_SET(int fd,fd_set *fdset);
   //关闭之前设置fdset中的bit位
   void FD_CLR(int fd,fd_set *fdset);
   //判断fd是否就绪
   void FD_ISSET(int fd,fd_set *fdset);

   实现原理：
   -1.传入的fd >> 5,保证其是在0-系统默认最多打开[2的系统位数]个文件描述符）
   -2.
*/

import (
	"fmt"
	"linux/ioTest"
	"syscall"
)

type Select struct {
}

func FD_ZERO(fdSet *syscall.FdSet) {
	for idx := range fdSet.Bits {
		fdSet.Bits[idx] = 0
	}
}

func FD_SET(fd int, fdSet *syscall.FdSet) {
	fdSet.Bits[(fd)>>5] |= 1 << ((fd) & 31)
}

func FD_CLR(fd int, fdSet *syscall.FdSet) {
	fdSet.Bits[(fd)>>5] &= ^(1 << (fd) & 31)
}

func FD_ISSET(fd int, fdSet *syscall.FdSet) bool {
	return (fdSet.Bits[(fd)>>5] & (1 << (fd) & 31)) != 0
}

// 启动server
func (s *Select) StartServer() error {
	//1.获取socketFd
	nSocketFd, err := ioTest.PreStartTcp()
	if err != nil {
		fmt.Println("Get socket error.", err)
		return err
	}
	fmt.Println(nSocketFd)
	//2.select 监听集合处理
	readFds := &syscall.FdSet{Bits: [16]int64{0}}
	maxFd := handleFdSet(nSocketFd, readFds)
	cliFds := make([]int, 100)
	//3.select 进行监听描述符
	for true {
		maxFd = handleFdSet(nSocketFd, readFds)
		//存放连接的client到readFds中
		for idx, val := range cliFds {
			FD_SET(val, readFds)
			if cliFds[idx] > maxFd {
				maxFd = cliFds[idx]
			}
		}

		//select 监控readFds集合,select的地一个参数是待测试的描述符个数+1
		_, err := syscall.Select(maxFd+1, readFds, nil, nil, nil)
		if err != nil {
			fmt.Printf("select error %s", err)
			return err
		}
		szBuff := make([]byte, 100)
		//判断client 发送数据
		for _, val := range cliFds {
			if FD_ISSET(val, readFds) {
				nRdSocketLen, err := syscall.Read(val, szBuff)
				if err != nil {
					fmt.Printf("[read] error %s", err)
					return err
				}
				_, err = syscall.Write(val, szBuff)
				if err != nil {
					fmt.Printf("[Write] error %s", err)
					return err
				}

				if nRdSocketLen > 0 {
					fmt.Printf("read data from client %s\n", szBuff)
				}
			}
		}

		//判断是否有client连接
		if FD_ISSET(nSocketFd, readFds) {
			fmt.Println(nSocketFd, readFds)
			cliFd, cliInfo, err := syscall.Accept(nSocketFd)
			if err != nil {
				fmt.Printf("[accept] error %s\n,cliFd %d\n", err, cliFd)
				return err
			}

			cliFds = append(cliFds, cliFd)
			fmt.Printf("clientInfo has connected:%v\n", cliInfo)
		}
	}
	err = syscall.Close(nSocketFd)
	if err != nil {
		fmt.Println("syscall.Close error")
	}
	return err
}

// 处理select要用的fdSet集合
func handleFdSet(nSocketFd int, readFds *syscall.FdSet) int {
	FD_ZERO(readFds)
	FD_SET(nSocketFd, readFds)
	return nSocketFd
}
