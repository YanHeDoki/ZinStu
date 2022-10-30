package ziface

//定义一个服务器接口

type IServer interface {
	//启动
	Start()
	//停止
	Stop()
	//运行
	Server()

	//路由功能：给当前的服务器注册一个路由方法。供客户端的链接处理使用
	AddRouter(msgid uint32, router IRouter)
	//返回连接资源管理器
	GetMgr() IConnManager

	//设置该Server的连接创建时Hook函数
	SetOnConnStart(func(IConnection))
	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func(IConnection))
	//调用连接OnConnStart Hook函数
	CallOnConnStart(conn IConnection)
	//调用连接OnConnStop Hook函数
	CallOnConnStop(conn IConnection)
}
