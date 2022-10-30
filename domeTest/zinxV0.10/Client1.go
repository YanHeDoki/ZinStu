package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/utils"
	"zinx/znet"
)

//模拟客户端
func main() {
	utils.ConfigInit()
	//1创建直接链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("net dial err:", err)
		return
	}
	defer conn.Close()
	//链接调用write方法写入数据
	for {
		dp := znet.NewDP()
		msg, err := dp.Pack(znet.NewMsgPackage(1, []byte("hello zinx-v0.10")))
		if err != nil {
			return
		}
		_, err = conn.Write(msg)
		if err != nil {
			return
		}
		headdata := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headdata)
		if err != nil {
			fmt.Println("read head error")
			break
		}
		msghead, err := dp.UnPack(headdata)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}
		if msghead.GetDataLen() > 0 {
			//msg 是有data数据的，需要再次读取data数据
			msg := msghead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}
			fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}
		time.Sleep(1 * time.Second)
	}
}
