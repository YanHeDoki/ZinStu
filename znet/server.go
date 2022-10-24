package znet

import (
	"errors"
	"fmt"
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

//暂且写死这个方法
func CallBack(conn *net.TCPConn, buf []byte, n int) error {
	wl, err := conn.Write(buf[:n])
	if err != nil {
		return errors.New("Read err")
	}
	fmt.Println("callback write len is ", wl)
	return nil
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
		var cid uint32
		cid = 0
		//阻塞的等待客户端的连接 处理客户端的链接业务（读写）
		for {
			conn, err := listen.AcceptTCP()
			if err != nil {
				fmt.Println("AcceptTCP err:", err)
				continue
			}
			//使用新的connection模块
			newConnection := NewConnection(conn, cid, CallBack)
			cid++
			newConnection.Start()
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
