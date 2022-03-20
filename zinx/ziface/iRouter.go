package ziface

type IRouter interface {
	// 处理conn业务之前的钩子方法Hook
	PreHandle(request IRequest)
	// 在处理conn业务的主方法hook
	Handle(request IRequest)
	// 在处理conn业务
	PostHandle(request IRequest)
}

