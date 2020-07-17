package ioEpoll

import (
	"fmt"
	"linux/ioTest"
	"syscall"
)

type Epoll struct {
}

func (e *Epoll) StartServer() error {
	//1.获取socketFd
	nSocketFd, err := ioTest.PreStartTcp()
	if err != nil {
		fmt.Println("Get socket error.", err)
		return err
	}

	//2-1.epoll_create 句柄,内核创建evnetpoll对象，文件系统一员
	epfd, err := syscall.EpollCreate(256)
	if err != nil {
		fmt.Printf("[epollCreate] has error %s\n", err)
		return err
	}
	//fmt.Printf("nSocketFd is %d\n",nSocketFd)
	//2-2.epoll 监听的事件结构体
	ev := &syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(nSocketFd),
	}
	fmt.Printf("nSocketFd is %d\n",nSocketFd)
	//2-3.注册 nSocketFd,添加到红黑树上
	err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, nSocketFd, ev)
	if err != nil {
		fmt.Printf("[epollCtl] has error %s\n", err)
		return err
	}
	var events []syscall.EpollEvent
	events = make([]syscall.EpollEvent, 20)

	for {
		//2-4. epoll wait 等待监听事件发生
		eventNums, err := syscall.EpollWait(epfd, events, 20)
		if err != nil {
			fmt.Printf("[epollwait] has error. %s\n", err)
			return err
		}
	//	fmt.Printf("eventNums is %d\n", eventNums)
		for index := 0; index < eventNums; index++ {
			//3-1.有新的client连接
			if events[index].Fd == int32(nSocketFd) {
				cliFd, cliInfo, err := syscall.Accept(nSocketFd)
				if err != nil {
					fmt.Printf("[accept] has error %s\n", err)
					return err
				}
				fmt.Printf("client %d has connected,clintInfo is %v\n", cliFd, cliInfo)
				//注册 cliFd
				ev.Fd = int32(cliFd)
				err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, cliFd, ev)
				if err != nil {
					fmt.Printf("[epollctl] has error %s", err)
					return err
				}
			} else {
				messageBuff := make([]byte, 1024)
				_, err := syscall.Read(int(events[index].Fd), messageBuff)
				if err != nil {
					fmt.Printf("[readFile] error %s", err)
					return err
				}
				fmt.Printf("Receive from client message:%s\n", messageBuff)
				_, err = syscall.Write(int(events[index].Fd), messageBuff)
				if err != nil {
					fmt.Printf("[writeToClient] error.%s\n", err)
					return err
				}
			}
		}
	}
	syscall.Close(epfd)
	syscall.Close(nSocketFd)
	return nil
}
