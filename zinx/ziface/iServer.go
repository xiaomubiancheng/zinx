package ziface

type IServer interface {
	Start()
	Stop()
	Serve()

	//路由
	AddRouter(msgID uint32,router IRouter)
}
