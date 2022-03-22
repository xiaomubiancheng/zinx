package ziface

type IMsgHandle interface {
	// 调度/执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)
	// 为消息添加具体的处理逻辑
	AddRouter(msgID uint32,router IRouter)

	// 启动Worker工作池
	StartWorkerPool()
	//
	SendMsgToTaskQueue(request IRequest)
}
