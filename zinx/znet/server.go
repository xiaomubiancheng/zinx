package znet

import (
	"fmt"
	"net"
	"zinxAll/zinx/ziface"
	"zinxAll/zinx/utils"
)

type Server struct{
	Name string
	IPVersion string
	IP string
	Port int
	//
	Router ziface.IRouter
}



func(s *Server)Start(){
	fmt.Printf("[start]Server Listenner at IP:%s,Port:%d,\n",s.IP,s.Port)

	go func(){
		//1.获取一个TCP的地址
		addr,err  := net.ResolveTCPAddr(s.IPVersion,fmt.Sprintf("%s:%d",s.IP,s.Port))
		if err!=nil{
			fmt.Println("resolve tcp addt error:",err)
		}
		// 2.监听服务器的地址
		listener,err := net.ListenTCP(s.IPVersion,addr)
		if err!=nil{
			fmt.Println("listen",s.IPVersion,"err",err)
			return
		}
		fmt.Println("start Zinx server succ,",s.Name,"succ,Listening...")
		var cid uint32
		cid = 0

		// 3.阻塞等待客户端链接,处理客户端链接业务
		for {
			//如果有客户端链接过来,阻塞会返回
			conn,err := listener.AcceptTCP()
			if err !=nil{
				fmt.Println("Accept err",err)
				continue
			}

			// 将处理新连接的业务方法和conn进行绑定
			dealConn := NewConnection(conn,cid,s.Router)
			cid++

			// 启动
			go dealConn.Start()
		}
	}()

}


func(s *Server)Stop(){

}

func(s *Server)Serve(){
	// 启动server
	s.Start()

	//TODO  else


	// 阻塞状态
	select{}
}

func(s *Server)AddRouter(router ziface.IRouter){
	s.Router = router
	fmt.Println("Add Router Succ!!")
}

func NewServer(name string)ziface.IServer{
	return &Server{
		Name:utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP: utils.GlobalObject.Host,
		Port:utils.GlobalObject.TcpPort,
		Router:nil,
	}
}