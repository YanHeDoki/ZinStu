package znet

import (
	"errors"
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

	//该链接处理的方法
	//消息的管理msgid对应的方法业务
	MsgHandler ziface.IMsgHandle
}

func (c *Connection) StartReader() {

	fmt.Println("Reader Server start ....")
	defer c.Stop()

	for {
		//buf := make([]byte, utils.GlobalConfig.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	//处理一下eof错误
		//	if err == io.EOF {
		//		err = nil
		//		continue
		//	}
		//	fmt.Println("Conn reader to buf err:", err)
		//	continue
		//}
		msghead := NewDP()
		headbuf := make([]byte, msghead.GetHeadLen())
		_, err := io.ReadFull(c.Conn, headbuf)
		if err != nil {
			fmt.Println("read in packhead err:", err)
			break
		}
		message, err := msghead.UnPack(headbuf)
		if err != nil {
			fmt.Println("UnPack err:", err)
			break
		}

		//根据datalen的参数再去读取一次
		var data []byte
		if message.GetDataLen() > 0 {
			data = make([]byte, message.GetDataLen())
			if _, err := io.ReadFull(c.GetTcpConnection(), data); err != nil {
				fmt.Println("read data err:", err)
				break
			}
		}
		message.SetData(data)

		//得到当前数据的Request 数据
		req := Request{
			conn: c,
			msg:  message,
		}

		//从路由中 找到注册绑定的Conn对应的router调用
		go c.MsgHandler.DoMsgHandler(&req)
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

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClosed {
		return errors.New("conn isClosed when send msg")
	}
	dp := NewDP()
	binaryMesg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		return err
	}
	_, err = c.Conn.Write(binaryMesg)
	if err != nil {
		return err
	}

	return nil
}

func NewConnection(conn *net.TCPConn, ConnID uint32, msghandler ziface.IMsgHandle) *Connection {
	return &Connection{
		Conn:       conn,
		ConnID:     ConnID,
		IsClosed:   false,
		ExitChan:   make(chan bool, 1),
		MsgHandler: msghandler,
	}

}
