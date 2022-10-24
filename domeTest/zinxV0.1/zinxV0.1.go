package main

import "zinx/znet"

func main() {
	server := znet.NewServer("test1")
	server.Server()
}
