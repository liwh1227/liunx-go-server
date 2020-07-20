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
*/

import (
	"fmt"
	"linux/ioTest"
	"syscall"
)

type Select struct {
}

// 启动server
func (s *Select) StartServer() error {
	//1.获取socketFd
	nSocketFd, err := ioTest.PreStartTcp()
	if err != nil {
		fmt.Println("Get socket error.", err)
		return err
	}
	//2.select 监听集合处理
	readFds := &syscall.FdSet{Bits: [16]int64{0}}
	allFds := &syscall.FdSet{Bits: [16]int64{0}}
	maxFd := handleFdSet(nSocketFd, allFds)
	fmt.Println("max fd is", maxFd)
	//存放要监听client的集合
	cliFds := make([]int, 0)
	maxi := -1
	//3.select 进行监听描述符
	for true {
		//重置fdset
		readFds = allFds
		//select 监控readFds集合,select的地一个参数是待测试的描述符个数+1
		fmt.Println("Before...", readFds)
		nReady, err := syscall.Select(maxFd+1, readFds, nil, nil, nil)
		fmt.Println("nReady is ", nReady)
		if err != nil {
			fmt.Printf("select error %s", err)
			return err
		}

		//判断是否有client连接
		if FD_ISSET(nSocketFd, readFds) {
			cliFd, cliInfo, err := syscall.Accept(nSocketFd)
			if err != nil {
				fmt.Printf("[accept] error %s\n,cliFd %d\n", err, cliFd)
				return err
			}

			length := 0
			var i int
			for i, v := range cliFds {
				if v == -1 {
					cliFds[i] = cliFd
					break
				}
				length++
			}
			if length == len(cliFds) {
				cliFds = append(cliFds, cliFd)
			}
			FD_SET(cliFd, allFds)
			fmt.Printf("clientInfo has connected:%v\n,cliFds is %v\n", cliInfo, cliFds)
			//将此时client的fd置为最大监控maxfd
			if cliFd > maxFd {
				maxFd = cliFd
			}
			//clients的最大下标
			if i > maxi {
				maxi = i
			}
			nReady--
			fmt.Println("nReady is ", nReady)
			if nReady <= 0 {
				continue
			}
		}

		for i, v := range cliFds {
			if nSocketFd == cliFds[i] {
				fmt.Println("dsdsadasdasdasdsadasd")
				continue
			}

			szBuff := make([]byte, 100)
			if FD_ISSET(v, readFds) {
				nRdSocketLen, err := syscall.Read(v, szBuff)
				if err != nil {
					fmt.Printf("[read] error %s", err)
					return err
				}
				if nRdSocketLen == 0 {
					fmt.Println("client has disconnected.")
					err := syscall.Close(v)
					if err != nil {
						return err
					}
					FD_CLR(v, allFds)
					cliFds[i] = -1
				} else if nRdSocketLen > 0 {
					fmt.Printf("read data from client:\n %s\n", szBuff)
				}
				nReady--
				fmt.Println("nReady line 114 is ", nReady)
				if nReady <= 0 {
					break
				}
			}
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
