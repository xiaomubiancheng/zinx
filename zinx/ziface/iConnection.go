package ziface

import "net"

type IConnection interface {
	Start()
	Stop()
	// 获取当前链接的绑定socket conn
	GetTCPConnection() *net.TCPConn
	GetConnID() uint32
	//远程客户端的 TCP状态
	RemoteAddr() net.Addr
	SendMsg(msgId uint32,data []byte)error

	// 设置链接属性
	SetProperty(key string,value interface{})
	// 获取链接属性
	GetProperty(key string)(interface{},error)
	//移除链接属性
	RemoveProperty(key string)
}

// 定义一个处理链接业务的方法
type HandleFunc func(*net.TCPConn,[]byte,int)error
