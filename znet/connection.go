package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

//当前链接模块
type Connection struct {
	//当前链接隶属于哪个Server
	TcpServer ziface.IServer
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
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	MsgBuffChan chan []byte //定义channel成员
	//该链接处理的方法
	//消息的管理msgid对应的方法业务
	MsgHandler ziface.IMsgHandle

	//链接属性
	property map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	res, ok := c.property[key]
	if !ok {
		return nil, errors.New("not found key for Connection Property ")
	} else {
		return res, nil
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
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

			//针对有缓冲channel需要些的数据处理
		case data, ok := <-c.MsgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				fmt.Println("MsgBuffChan is Closed")
				break
			}
		//用于退出
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
	//调用开发者设置的启动前的钩子函数
	c.TcpServer.CallOnConnStart(c)
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
	//在销毁连接之前执行开发者的函数
	c.TcpServer.CallOnConnStop(c)
	c.Conn.Close()

	//结束 通知chan已经结束了连接
	c.ExitChan <- true
	close(c.ExitChan)
	//将当前链接从connmgr中销毁
	c.TcpServer.GetMgr().Remove(c)
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
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		return err
	}
	c.MsgChan <- msg
	return nil
}

//带缓冲发送消息
func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.IsClosed == true {
		return errors.New("Connection closed when send buff msg")
	}
	//将data封包，并且发送
	dp := NewDP()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}

	//写回客户端
	c.MsgBuffChan <- msg
	return nil
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, ConnID uint32, msghandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:        conn,
		ConnID:      ConnID,
		IsClosed:    false,
		ExitChan:    make(chan bool, 1),
		MsgHandler:  msghandler,
		MsgChan:     make(chan []byte),
		MsgBuffChan: make(chan []byte, utils.GlobalConfig.MaxWorkerTaskLen), //不要忘记初始化
		TcpServer:   server,
		property:    make(map[string]interface{}),
	}
	c.TcpServer.GetMgr().Add(c)
	return c
}
