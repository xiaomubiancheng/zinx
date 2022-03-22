package ziface

type IServer interface {
	Start()
	Stop()
	Serve()

	//路由
	AddRouter(uint32,IRouter)
	GetConnMgr() IConnManager

	// 注册OnConnStart 钩子函数的方法
	SetOnConnStart(func(connection IConnection))
	//
	SetOnConnStop(func(connection IConnection))
	// 调用OnConnStart钩子函数的方法
	CallOnConnStart(connection IConnection)
	CallOnConnStop(connection IConnection)

}
