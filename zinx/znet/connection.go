package znet

import (
	"fmt"
	"net"
	"zinxAll/zinx/utils"
	"zinxAll/zinx/ziface"
)

type Connection struct{
	Conn *net.TCPConn
	ConnID uint32
	isClosed bool

	// 告知当前链接已经退出的/停止 channel
	ExitChan chan bool

	//路由
	Router ziface.IRouter
}

func NewConnection(conn *net.TCPConn,connID uint32,router ziface.IRouter)*Connection{
	c := &Connection{
		Conn:      conn,
		isClosed:  false,
		ConnID:    connID,
		Router: router,
		ExitChan:  make(chan bool,1),
	}
	return c
}


func (c *Connection)StartReader(){
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID=",c.ConnID,"Reader is exit,remote addr id",c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buf中,最大512字节
		buf := make([]byte,utils.GlobalObject.MaxPackageSize)
		_,err := c.Conn.Read(buf)
		if err!=nil{
			fmt.Println("recv buf err",err)
			continue
		}

		//得到当前conn数据的Request请求数据
		req := Request{
			conn:c,
			data: buf,
		}

		// 从路由中，找到注册绑定的Conn对应的router调用
		go func(req ziface.IRequest){
			c.Router.PreHandle(req)
			c.Router.Handle(req)
			c.Router.PostHandle(req)
		}(&req)

	}
}

func (c *Connection)Start(){
	fmt.Println("Conn Start() ... ConnID=",c.ConnID)
	//启动从当前链接的读数据的业务
	go c.StartReader()
	//TODO
}
func (c *Connection)Stop(){
	fmt.Println("Conn Stop().. ConnID=", c.ConnID)

	//如果当前链接已经关闭
	if c.isClosed == true{
		return
	}
	c.isClosed = true

	//关闭socket链接
	c.Conn.Close()

	//回收资源
	close(c.ExitChan)
}

func (c *Connection)GetTCPConnection() *net.TCPConn{

	return c.Conn

}
func (c *Connection)GetConnID() uint32{
	return c.ConnID

}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	return nil
}




