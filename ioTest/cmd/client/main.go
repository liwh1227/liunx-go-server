package main

import (
	"fmt"
	"net"
)

func main() {
	i := 0
	maxRoutine := 12800
	for i < maxRoutine {
		go func() {
			client()
		}()
		i++
		//fmt.Println(i)
		//time.Sleep(100 * time.Microsecond)
	}
	select {}
}
func client() {
	conn, err := net.Dial("tcp", "localhost:8080")
	defer conn.Close()
	if err != nil {
		fmt.Printf("connect failed, err : %v\n", err.Error())
		return
	}

	for {
		//将line发给服务器
		n, err := conn.Write([]byte("hello world"))
		if err != nil {
			fmt.Println("conn.Write err=", err)
		}
		fmt.Printf("发送了%d个字节\n", n)

		//创建切片
		buf := make([]byte, 1024)

		//2 如果没有writer发送就一直阻塞在这
		n, err = conn.Read(buf)
		if err != nil {
			fmt.Println("服务器read err=", err) //出错退出
			return
		}
		//3. 显示读取内容到终端
		fmt.Print(string(buf[:n]))
	}
}
