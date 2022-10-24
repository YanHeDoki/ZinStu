package ziface

//路由接口
//路由里的接口都是IRequest

type IRouter interface {
	//处理conn业务之前的方法Hook
	PreHandle(request IRequest)
	//处理conn业务的方法
	Handle(request IRequest)
	//处理conn业务之后的方法Hook
	AfterHandle(request IRequest)
}
