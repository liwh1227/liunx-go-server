package single

/*
	1.单进程与客户端通信
    2.缺点：若同时有多个client连接后，使用该模型，由于进程陷入内核调用中等待，server无法同时处理多个客户端信息,
*/
import (
	"bufio"
	"fmt"
	"net"
)

type SingleThread struct {
}

func handleConnection(conn net.Conn) error {
	//read也会阻塞
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Printf("client send msg: %s\n", status)
	_, err = conn.Write([]byte("hello client" + "\n"))
	if err != nil {
		fmt.Println("conn.Write err=", err)
	}
	return err
}

func (s *SingleThread) StartServer() error {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("[listen] has error %s\n", err)
		return err
	}
	for true {
		//accept阻塞
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		fmt.Printf("[client] has connected,%v\n", conn)
		err = handleConnection(conn)
		if err != nil {
			fmt.Printf("[handleConnection] has error.%s\n", err)
		}
	}
	return err
}
