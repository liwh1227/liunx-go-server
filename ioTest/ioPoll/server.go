package io_poll

import (
	"fmt"
	"golang.org/x/sys/unix"
	"linux/ioTest"
	"time"
)

type Poll struct {
}

func (p *Poll) StartServer() error {
	//1.获取socketFd
	var cliNum int
	var t1 time.Time
	listenFd, err := ioTest.PreStartTcp()
	if err != nil {
		fmt.Println("Get socket error.", err)
		return err
	}
	clients := make([]unix.PollFd, 15000)
	clients[0].Fd = int32(listenFd)
	clients[0].Events = unix.POLLIN

	for i := 1; i < len(clients); i++ {
		clients[i].Fd = -1
	}
	maxi := 0
	for {
		nReady, err := unix.Poll(clients, -1)
		if cliNum == 0 {
			t1 = time.Now()
		}
		if err != nil {
			fmt.Printf("[poll] error %s\n", err)
			return err
		}

		if clients[0].Revents&unix.POLLIN > 0 {
			connfd, _, err := unix.Accept(int(clients[0].Fd))
			if err != nil {
				fmt.Printf("[accept] has error %s\n", err)
				return err
			}
			//fmt.Printf("[client] %v has connected\n", clientInfo)
			cliNum++
			if cliNum == 12800 {
				et := time.Since(t1)
				fmt.Println(et)
			}
			fmt.Println(cliNum)
			var i int
			for i = range clients {
				if clients[i].Fd < 0 {
					clients[i].Fd = int32(connfd)
					break
				}
			}
			if i > 15000 {
				fmt.Println("Too many clients")
				return err
			}

			clients[i].Events = unix.POLLIN

			if i > maxi {
				maxi = i
			}

			nReady--
			if nReady <= 0 {
				continue
			}
		}

		for i := 1; i <= maxi; i++ {
			sockFd := clients[i].Fd
			if sockFd < 0 {
				continue
			}
			messageBuffer := make([]byte, 100)
			if clients[i].Revents&(unix.POLLIN|unix.POLLERR) > 0 {
				nRead, err := unix.Read(int(sockFd), messageBuffer)
				//fmt.Println(messageBuffer)
				if err != nil {
					fmt.Println("[read] error", err)
					return err
				} else if nRead == 0 {
					fmt.Println("[client] has disconnected with server")
					err = unix.Close(int(sockFd))
					if err != nil {
						fmt.Println("close sockFd err", err)
						return err
					}
					clients[i].Fd = -1
				} else if nRead > 0 {
					fmt.Printf("Receive from client message:%s\n", messageBuffer)
					//fmt.Printf("client send message:---%s\n", messageBuffer)
				}

				nReady--
				if nReady <= 0 {
					break
				}
			}
		}
	}
	return err
}
