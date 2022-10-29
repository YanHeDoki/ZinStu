package znet

import (
	"zinx/ziface"
)

//实现router时先嵌入这个基类，然后根据需求重写其中的方法就是了
type BaseRouter struct {
}

//之所以把这些方法都留空是因为只需要实现主业务方法时候我们只需要继承baseRouter然后实现主业务方法就好了
func (b *BaseRouter) PreHandle(request ziface.IRequest) {

}

func (b *BaseRouter) Handle(request ziface.IRequest) {

}

func (b *BaseRouter) AfterHandle(request ziface.IRequest) {

}
