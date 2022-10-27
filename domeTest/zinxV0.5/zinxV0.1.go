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

func (p *PingRouter) PreHandle(request ziface.IRequest) {
}

func (p *PingRouter) AfterHandle(request ziface.IRequest) {
}

// test handleRouter
func (p *PingRouter) Handle(req ziface.IRequest) {

	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from data=", string(req.GetData()))

	if err := req.GetConnection().SendMsg(1, []byte("ping...ping...ping")); err != nil {
		fmt.Println(err)
	}
}

func main() {
	fmt.Println("test Zinx-V0.5")
	//加载配置文件
	utils.ConfigInit()
	fmt.Println(utils.GlobalConfig)
	server := znet.NewServer()
	p := &PingRouter{}
	server.AddRouter(p)
	server.Server()

}
