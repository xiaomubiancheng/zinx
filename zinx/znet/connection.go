package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
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
}

func NewConnection(conn *net.TCPConn,connID uint32,msgHandler ziface.IMsgHandle)*Connection{
	c := &Connection{
		Conn:      conn,
		isClosed:  false,
		ConnID:    connID,
		MsgHandler: msgHandler,
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

		// 从路由中，找到注册绑定的Conn对应的router调用
		go c.MsgHandler.DoMsgHandler(&req)

	}
}





func (c *Connection)Start(){
	fmt.Println("Conn Start() ... ConnID=",c.ConnID)
	//启动从当前链接的读数据的业务
	go c.StartReader()
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
	if _,err := c.Conn.Write(binaryMsg);err!=nil{
		fmt.Println("Write msg id",msgId,"error:",err)
		return errors.New(" conn Write error ")
	}

	return nil
}




