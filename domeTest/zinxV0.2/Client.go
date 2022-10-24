package main

import (
	"fmt"
	"net"
)

//模拟客户端
func main() {

	//1创建直接链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("net dial err:", err)
		return
	}
	defer conn.Close()
	//链接调用write方法写入数据
	conn.Write([]byte("hello world2"))
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	fmt.Println(string(buf), n)

}
