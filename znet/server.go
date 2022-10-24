package znet

import (
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

type Server struct {
	//服务器名称
	Name string
	//IP版本 IPv4 or other
	IPVersion string
	//服务器ip
	IP string
	//服务器监听的端口
	Port int
}

func (s *Server) Start() {

	//日志，以后应该用日志来处理
	fmt.Printf("[start]Server Listenner at IP %s,Port %d ,is staring", s.IP, s.Port)

	//由server方法来阻塞所以异步处理

	go func() {
		//获取一个Tcp的Addr地址
		resolveIPAddr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("Start ServerErr err:", err)
			return
		}
		//监听服务器的地址
		listen, err := net.ListenTCP(s.IPVersion, resolveIPAddr)
		if err != nil {
			fmt.Println("ListenIPErr err:", err)
			return
		}

		fmt.Println("Start Zinx Server success", s.Name, "success Listening...")

		//阻塞的等待客户端的连接 处理客户端的链接业务（读写）
		for {
			conn, err := listen.AcceptTCP()
			if err != nil {
				fmt.Println("AcceptTCP err:", err)
				continue
			}
			//我们这里暂时做一个最大512字节的回显服务

			//用go来异步处理
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						if err == io.EOF {
							err = nil
							continue
						}
						fmt.Println("Read err:", err)
						continue
					}
					_, err = conn.Write(buf[:cnt])
					if err != nil {
						fmt.Println("Write err:", err)
						continue
					}
				}
			}()
		}
	}()

}

func (s *Server) Stop() {

	//断开服务器，将一些服务器的资源链接释放

}

func (s *Server) Server() {

	//启动服务器
	s.Start()

	//TODO 留空位可以给以后操作空间
	//阻塞 否则主Go退出， listenner的go将会退出
	//select {}
	for {

	}
}

//初始化server服务器方法
func NewServer(name string) ziface.IServer {
	return &Server{ //报错不能返回这个类型
		Name:      name,
		IPVersion: "tcp4",
		Port:      8999,
	}
}