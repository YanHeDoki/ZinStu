package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)

type MsgHandle struct {
	Apis map[uint32]ziface.IRouter
}

//构造方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

func (m *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handle, ok := m.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("not find Router In Apis")
		return
	}
	handle.PreHandle(request)
	handle.Handle(request)
	handle.AfterHandle(request)

}

func (m *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := m.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	//2 添加msg与api的绑定关系
	m.Apis[msgId] = router
	fmt.Println("Add api msgId = ", msgId)
}
