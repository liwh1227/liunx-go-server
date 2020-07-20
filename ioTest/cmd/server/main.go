package main

import (
	"fmt"
	io_poll "linux/ioTest/ioPoll"
)

func main() {
	//eSer := &ioEpoll.Epoll{}
	//err := eSer.StartServer()
	//if err != nil {
	//	fmt.Println("start server err")
	//}

	/*sig := &single.SingleThread{}
	err := sig.StartServer()
	if err != nil {
		fmt.Printf("[StartServer] has error.%s\n",err)
	}*/

	/*selectSer := ioSelect.Select{}
	err := selectSer.StartServer()
	if err != nil {
		fmt.Printf("[select io] has error %s\n",selectSer)
	}*/
	pollSer := &io_poll.Poll{}
	err := pollSer.StartServer()
	if err != nil {
		fmt.Printf("err is %s\n", err)
	}
}
