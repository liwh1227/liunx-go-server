package ioEpoll

/*	参照官方文档：https://www.man7.org/linux/man-pages/man2/epoll_ctl.2.html
    desc:
    This system call is used to add, modify, or remove entries in the
    interest list of the epoll(7) instance referred to by the file
    descriptor epfd.  It requests that the operation op be performed for
    the target file descriptor, fd.
	epoll 是linux下，常见的io复用模型之一，其通过底层双链表和红黑树的结构，能够高效的监听文件描述符fd，并返回结果。
	1.epoll_create:int epoll_create(int size);
	*注：epoll_create() creates a new epoll(7) instance.  Since Linux 2.6.8, the size argument is ignored.
	返回的int即内核分配的epollFd描述符，为了后续与监听描述符建立映射使用。
    2.epoll_ctl:int epoll_ctl(int epfd, int op, int fd, struct epoll_event *event);
    epfd:文件句柄即描述符;
    int op ：EPOLL_CTL_ADD、EPOLL_CTL_MOD、EPOLL_CTL_DEL 添加、修改、删除fd参数
	int fd：要监听的文件描述符
    *event:与fd建立联系的结构体，其中events表示描述符fd的events事件，重点关心，et和lt两种工作模式
*/
import (
	"fmt"
	"golang.org/x/sys/unix"
	"linux/ioTest"
	//"unix"
	"time"
)

type Epoll struct {
}

func (e *Epoll) StartServer() error {
	var cliNum int
	var t1 time.Time
	//1.获取socketFd
	nSocketFd, err := ioTest.PreStartTcp()
	if err != nil {
		fmt.Println("Get socket error.", err)
		return err
	}

	//2.获取epfd句柄，并注册socketFd描述符
	epfd, err := handleEpoll(nSocketFd, unix.EPOLLIN)
	if err != nil {
		return err
	}
	var events []unix.EpollEvent
	events = make([]unix.EpollEvent, 13000)
	for i := 1; i < len(events); i++ {
		events[i].Fd = -1
	}
	maxi := 0
	for true {
		//2-4. epoll wait 等待监听事件发生
		nReady, err := unix.EpollWait(epfd, events, 13000)
		//fmt.Println(nReady)
		if err != nil {
			fmt.Printf("[epollwait] has error. %s\n", err)
			return err
		}
		if cliNum == 0 {
			t1 = time.Now()
		}
		for i := 0; i < nReady; i++ {
			//3-1.有新的client连接
			if events[i].Fd == int32(nSocketFd) {
				connfd, _, err := unix.Accept(nSocketFd)
				if err != nil {
					fmt.Printf("[accept] has error %s\n", err)
					return err
				}
				//fmt.Printf("client %d has connected,clintInfo is %v\n", connfd, cliInfo)
				cliNum++
				if cliNum == 12800 {
					et := time.Since(t1)
					fmt.Println(et)
				}
				var i int
				for i = range events {
					if events[i].Fd < 0 {
						events[i].Fd = int32(connfd)
						break
					}
				}
				if i > 15000 {
					fmt.Println("Too many clients")
					return err
				}

				events[i].Events = unix.EPOLLIN | unix.EPOLLET

				if i > maxi {
					maxi = i
				}
				//注册 cliFd
				err = unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, connfd, &events[i])
				if err != nil {
					fmt.Printf("[epollctl] has error %s", err)
					return err
				}
				nReady--
				if nReady <= 0 {
					continue
				}
				//break
				//输入缓冲区输入
			} else {
				sockFd := events[i].Fd
				if sockFd < 0 {
					fmt.Println(sockFd)
					continue
				}
				messageBuff := make([]byte, 100)
				nRead, err := unix.Read(int(sockFd), messageBuff)
				if err != nil {
					fmt.Printf("[readFile] error %s", err)
					return err
				}
				if nRead == 0 {
					fmt.Printf("client has disconnected...\n")
					err = unix.Close(int(events[i].Fd))
					if err != nil {
						fmt.Printf("sys close line 80")
						return err
					}
					events[i].Fd = -1
				} else if nRead > 0 {
					fmt.Printf("Receive from client message:%s\n", messageBuff)
					//_,err := unix.Write(int(sockFd),messageBuff)
					//if err != nil {
					//	return err
					//}
				}
				nReady--
				if nReady <= 0 {
					break
				}
			}
		}
	}
	err = unix.Close(epfd)
	if err != nil {
		fmt.Printf("[close] epfd error %s\n", err)
		return err
	}
	err = unix.Close(nSocketFd)
	if err != nil {
		fmt.Printf("[close] nSocketFd error %s\n", err)
		return err
	}
	return nil
}

func handleEpoll(nSocketFd int, eventsMode uint32) (int, error) {
	//1.epoll_create 句柄,内核创建evnetpoll对象，文件系统一员，即epfd也占用一个文件描述符
	epfd, err := unix.EpollCreate(1)
	if err != nil {
		fmt.Printf("[epollCreate] has error %s\n", err)
		return epfd, err
	}
	//2.epoll 监听的事件结构体
	ev := &unix.EpollEvent{
		Events: eventsMode,
		Fd:     int32(nSocketFd),
	}
	//3.注册 nSocketFd,添加到红黑树上
	err = unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, nSocketFd, ev)
	if err != nil {
		fmt.Printf("[epollCtl] has error %s\n", err)
		return epfd, err
	}
	return epfd, err
}
