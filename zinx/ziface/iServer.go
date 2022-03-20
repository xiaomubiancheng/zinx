package ziface

type IServer interface {
	Start()
	Stop()
	Serve()

	//路由
	AddRouter(router IRouter)
}
