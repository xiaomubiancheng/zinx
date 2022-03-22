package znet

import (
	"fmt"
	"strconv"
	"zinxAll/zinx/utils"
	"zinxAll/zinx/ziface"
)

type MsgHandle struct{
	// 存放每个MsgID 所对应的处理方法
	Apis map[uint32]ziface.IRouter
	//  负载Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作Worker池的worker数量
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle{
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
		TaskQueue: make([]chan ziface.IRequest,utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize:utils.GlobalObject.WorkerPoolSize,
	}
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

// 启动一个worker工作池(开启只发生一次)
func (mh *MsgHandle)StartWorkerPool(){
	// 根源workerPoolSize 分别开启Worker
	for i:=0;i<int(mh.WorkerPoolSize);i++{
		mh.TaskQueue[i] = make(chan ziface.IRequest,utils.GlobalObject.MaxWorkerTaskLen )
		go mh.StartOneWorker(i,mh.TaskQueue[i])
	}
}


// 启动一个Worker工作流程
func (mh *MsgHandle)StartOneWorker(workerID int,taskQueue chan ziface.IRequest){
	fmt.Println("Worker ID =",workerID,"is started...")
	// 不断的阻塞等待对应消息队列的消息
	for {
		select{
			// 如果有消息过来,出列的就是一个客户端的Request,执行当前Request
			case request := <-taskQueue:
				mh.DoMsgHandler(request)
		}
	}
}


//// 将消息交给TaskQueue， 由worker进行处理
func (mh *MsgHandle)SendMsgToTaskQueue(request ziface.IRequest){
	// 1.将消息平均分配给不通过的worker
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=",request.GetConnection().GetConnID(),"request MsgID=",request.GetMsgID(),"to WorkerID=",workerID)

	// 2.将消息发给对应的worker的TaskQueue即可
	mh.TaskQueue[workerID]<-request
}