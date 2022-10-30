package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

//链接管理模块

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection, 10000),
	}

}

func (c *ConnManager) Add(connection ziface.IConnection) {
	//保护共享资源 加锁
	c.connLock.Lock()
	defer c.connLock.Unlock()
	c.connections[connection.GetConnID()] = connection
	fmt.Println("ADD conn to manager success")
}

func (c *ConnManager) Remove(connection ziface.IConnection) {
	//保护共享资源 加锁
	c.connLock.Lock()
	defer c.connLock.Unlock()
	delete(c.connections, connection.GetConnID())
	fmt.Println("Remove conn to manager success")
}

func (c *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	//保护共享资源 加锁
	c.connLock.RLock()
	defer c.connLock.RUnlock()
	conn, ok := c.connections[connId]
	if !ok {
		return nil, errors.New("not find Conn in Server")
	}
	fmt.Println("Get conn success")
	return conn, nil
}

func (c *ConnManager) Len() int {
	return len(c.connections)
}

func (c *ConnManager) ClearConn() {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	for connid, conn := range c.connections {
		//停止这个连接的资源
		conn.Stop()
		delete(c.connections, connid)
	}
	fmt.Println("clear ConnManagerMap success")
}
