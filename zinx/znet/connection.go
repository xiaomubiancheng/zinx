package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
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
	MsgHandler ziface.IMsgHandle

	// 无缓冲的管道，用于读、写Goroutine之间的消息通信
	msgChan chan []byte

	//隶属于哪个server
	TcpServer ziface.IServer

	//链接属性集合
	property map[string]interface{}
	// 保护链接属性的锁
	propertyLock sync.RWMutex
}

func NewConnection(server ziface.IServer,conn *net.TCPConn,connID uint32,msgHandler ziface.IMsgHandle)*Connection{
	c := &Connection{
		Conn:      conn,
		isClosed:  false,
		ConnID:    connID,
		MsgHandler: msgHandler,
		ExitChan:  make(chan bool,1),
		msgChan: make(chan []byte),
		TcpServer: server,
		property: make(map[string]interface{}),
	}
	// 将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}


func (c *Connection)StartReader(){
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID=",c.ConnID,"Reader is exit,remote addr id",c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buf中,最大512字节
		//buf := make([]byte,utils.GlobalObject.MaxPackageSize)
		//_,err := c.Conn.Read(buf)
		//if err!=nil{
		//	fmt.Println("recv buf err",err)
		//	continue
		//}

		// 创建一个拆包解包对象
		dp := NewDataPack()
		// 读取客户端的Msg Head 二进制流8个字节
		headData := make([]byte,dp.GetHeadLen())
		if _,err :=io.ReadFull(c.GetTCPConnection(),headData);err!=nil{
			fmt.Println("read msg head error",err)
			break
		}
		// 拆包，得到msgID和msgDatalen 放在msg消息中
		msg ,err := dp.Unpack(headData)
		if err!=nil{
			fmt.Println("unpack error",err)
			break
		}
		// 根据dataLen 再次读取Data, 放在msg.Data中
		var data []byte
		if msg.GetMsgLen()>0{
			data = make([]byte,msg.GetMsgLen())
			if _,err := io.ReadFull(c.GetTCPConnection(),data);err!=nil{
				fmt.Println("read msg data error",err)
				break
			}
		}
		msg.SetData(data)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn:c,
			msg:msg,
		}

		if utils.GlobalObject.WorkerPoolSize>0{
			c.MsgHandler.SendMsgToTaskQueue(&req)

		}else{
			// 从路由中，找到注册绑定的Conn对应的router调用
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}


// 写消息的Goroutine, 发送给客户端
func (c *Connection)StartWriter(){
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(),"[conn Writer exit!]")

	for{
		select{
			case data:=<-c.msgChan:
				//有数据要写给客户端
				if _,err := c.Conn.Write(data);err!=nil{
					fmt.Println("Send data error,", err)
					return
				}
			case <-c.ExitChan:
				//代表Reader已经退出，此时Writer也要退出
				return
		}
	}
}



func (c *Connection)Start(){
	fmt.Println("Conn Start() ... ConnID=",c.ConnID)
	//启动从当前链接的读数据的业务
	go c.StartReader()

	go c.StartWriter()

	// 执行开发者的hook
	c.TcpServer.CallOnConnStart(c)
}
func (c *Connection)Stop(){
	fmt.Println("Conn Stop().. ConnID=", c.ConnID)

	//如果当前链接已经关闭
	if c.isClosed == true{
		return
	}
	c.isClosed = true

	//调用开发者注册的 销毁链接之前 需要执行的hook
	c.TcpServer.CallOnConnStop(c)

	//关闭socket链接
	c.Conn.Close()

	// 告知Writer关闭
	c.ExitChan <- true

	// 将当前连接从connMgr中摘除掉
	c.TcpServer.GetConnMgr().Remove(c)

	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
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

// 提供一个SendMsg方法 将我们要发送给客户端的数据，先进行封包，再发送
func(c *Connection)SendMsg(msgId uint32,data []byte)error{
	if c.isClosed == true{
		return errors.New(" Connection closed when send msg ")
	}

	// 将data 进行封包 MsgDataLen|MsgID|Data
	dq := NewDataPack()
	binaryMsg,err := dq.Pack(NewMsgPackage(msgId,data))
	if err!=nil{
		fmt.Println("Pack error msg id=",msgId)
		return errors.New("Pack error msg ")
	}
	//将数据发送给客户端
	c.msgChan <-binaryMsg

	return nil
}


func(c *Connection)SetProperty(key string,value interface{}){
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}
// 获取链接属性
func(c *Connection)GetProperty(key string)(interface{},error){
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value,ok:= c.property[key];ok{
		return value,nil
	}else{
		return nil,errors.New("no property found")
	}
}
//移除链接属性
func(c *Connection)RemoveProperty(key string){
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property,key)
}




