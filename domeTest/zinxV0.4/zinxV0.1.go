package main

import (
	"fmt"
	"zinx/utils"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	BR znet.BaseRouter
}

// test PreRouter
func (p *PingRouter) PreHandle(req ziface.IRequest) {
	fmt.Println("PreHandle test")
	req.GetConnection().GetTcpConnection().Write([]byte("before in handle"))
}

// test handleRouter
func (p *PingRouter) Handle(req ziface.IRequest) {
	req.GetConnection().GetTcpConnection().Write([]byte("ping handle"))
}

// test atfRouter
func (p *PingRouter) AfterHandle(req ziface.IRequest) {
	fmt.Println("AtfHandle test")
	req.GetConnection().GetTcpConnection().Write([]byte("aft in handle"))
}

func main() {
	fmt.Println("test Zinx-V0.4")
	//加载配置文件
	utils.ConfigInit()
	server := znet.NewServer()
	p := &PingRouter{}
	server.AddRouter(p)
	server.Server()

}
