package main

import "zinx/znet"

func main() {
	server := znet.NewServer("test2")
	server.Server()
}
