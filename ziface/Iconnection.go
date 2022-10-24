package ziface

import "net"

type IConnection interface {

	//开始链接方法
	Start()
	//停止链接方法
	Stop()
	//获取当前链接绑定的socket
	GetTcpConnection() *net.TCPConn
	//获取当前链接模块的链接ID
	GetConnID() uint32
	//获取远程客户端的Tcp状态 ip port
	RemoteAddr() net.Addr
	//发送数据，将数据发送给远程客户端
	Send(data []byte) error
}

//处理函数类型
type HandleFunc func(*net.TCPConn, []byte, int) error
