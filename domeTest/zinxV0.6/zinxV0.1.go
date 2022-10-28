package main

import (
	"fmt"
	"zinx/utils"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

// test handleRouter
func (p *PingRouter) Handle(req ziface.IRequest) {

	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from data=", req.GetMsgId(), string(req.GetData()))

	if err := req.GetConnection().SendMsg(0, []byte("ping...ping...ping")); err != nil {
		fmt.Println(err)
	}
}

//HelloZinxRouter Handle
type HelloZinxRouter struct {
	znet.BaseRouter
}

func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgId(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("Hello Zinx Router V0.6"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	fmt.Println("test Zinx-V0.6")
	//加载配置文件
	utils.ConfigInit()
	fmt.Println(utils.GlobalConfig)
	server := znet.NewServer()
	server.AddRouter(0, &PingRouter{})
	server.AddRouter(1, &HelloZinxRouter{})
	server.Server()

}
