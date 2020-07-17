package main

import (
	"fmt"
	"linux/ioTest/ioEpoll"

	//	"linux/ioTest/ioEpoll"
	//"linux/ioTest/ioSelect"
)

func main() {
	eSer := &ioEpoll.Epoll{}
	err := eSer.StartServer()
	if err != nil {
		fmt.Println("start server err")
	}
}
