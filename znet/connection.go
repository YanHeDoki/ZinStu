package znet

import (
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

//当前链接模块
type Connection struct {

	//当前链接的 socket tcp 套接字
	Conn *net.TCPConn
	//链接的ID
	ConnID uint32

	//当前链接的状态
	IsClosed bool

	//当前链接所绑定的处理业务方法的API
	HandleApi ziface.HandleFunc

	//告知当前链接已经退出的channel
	ExitChan chan bool
}

func (c *Connection) StartReader() {

	fmt.Println("Reader Server start ....")
	defer c.Stop()

	for {
		buf := make([]byte, 512)

		n, err := c.Conn.Read(buf)
		if err != nil {
			//处理一下eof错误
			if err == io.EOF {
				err = nil
				continue
			}
			fmt.Println("Conn reader to buf err:", err)
			continue
		}

		//绑定到具体业务处理上
		if err = c.HandleApi(c.Conn, buf, n); err != nil {
			fmt.Println("to handleApi err:", err)
			break
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("conn starting...ConnID=", c.ConnID)
	//启动当前链接的读数据业务
	go c.StartReader()
	//todo 启动当前链接的读数据业务

}

func (c *Connection) Stop() {
	fmt.Println("conn closeing...")

	if c.IsClosed {
		return
	} else {
		c.IsClosed = true
		c.Conn.Close()
		close(c.ExitChan)
	}
	fmt.Println("conn close suucess")
}

func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	return nil
}

func NewConnection(conn *net.TCPConn, ConnID uint32, handleapi ziface.HandleFunc) *Connection {
	return &Connection{
		Conn:      conn,
		ConnID:    ConnID,
		HandleApi: handleapi,
		IsClosed:  false,
		ExitChan:  make(chan bool, 1),
	}

}
