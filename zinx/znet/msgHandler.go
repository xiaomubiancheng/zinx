package znet

import (
	"fmt"
	"strconv"
	"zinxAll/zinx/ziface"
)

type MsgHandle struct{
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandle() *MsgHandle{
	return &MsgHandle{Apis: make(map[uint32]ziface.IRouter)}
}

func(mh *MsgHandle)DoMsgHandler(request ziface.IRequest){
	// 从Request中找到msgID
	handler,ok:= mh.Apis[request.GetMsgID()]
	if !ok{
		fmt.Println("api msgID=",request.GetMsgID(),"is NOT FOUND!Need Register")
	}
	// 根据MsgID调度对应router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func(mh *MsgHandle)AddRouter(msgID uint32,router ziface.IRouter){
	// 判断 当前msg绑定的API处理方法是否已经存在
	if _,ok:= mh.Apis[msgID];ok{
		// id已经注册了
		panic(any("repeat api,msgID="+strconv.Itoa(int(msgID))))
	}

	// 2.添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID=", msgID, "succ!")
}