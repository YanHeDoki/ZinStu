package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/utils"
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
	//无缓冲管道，用于读、写两个goroutine之间的消息通信
	MsgChan chan []byte
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

		//已经设置开启了工作池
		if utils.GlobalConfig.WorkerPoolSize > 0 {
			//发送消息到消息队列由工作池来处理
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//从路由中 找到注册绑定的Conn对应的router调用
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

/*
	写消息Goroutine， 用户将数据发送给客户端
*/
func (c *Connection) StartWrite() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data := <-c.MsgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error: ", err, " Conn Writer exit")
				return
			}
		case <-c.ExitChan:
			//conn已经关闭
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("conn starting...ConnID=", c.ConnID)
	//启动当前链接的读数据业务
	go c.StartReader()
	// 启动当前链接的读数据业务
	go c.StartWrite()

	for {
		select {
		case <-c.ExitChan:
			//得到退出消息，不再阻塞
			return
		}
	}
}

func (c *Connection) Stop() {
	fmt.Println("conn closeing...")

	if c.IsClosed {
		return
	}

	c.IsClosed = true
	c.Conn.Close()

	//结束 通知chan已经结束了连接
	c.ExitChan <- true
	close(c.ExitChan)
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
	c.MsgChan <- binaryMesg
	return nil
}

func NewConnection(conn *net.TCPConn, ConnID uint32, msghandler ziface.IMsgHandle) *Connection {
	return &Connection{
		Conn:       conn,
		ConnID:     ConnID,
		IsClosed:   false,
		ExitChan:   make(chan bool, 1),
		MsgHandler: msghandler,
		MsgChan:    make(chan []byte),
	}

}
