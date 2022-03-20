package znet

import (
	"net"
	"fmt"
	"zinxAll/zinx/ziface"
)

type Server struct{
	Name string
	IPVersion string
	IP string
	Port int
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
		// 3.阻塞等待客户端链接,处理客户端链接业务
		for {
			//如果有客户端链接过来,阻塞会返回
			conn,err := listener.AcceptTCP()
			if err !=nil{
				fmt.Println("Accept err",err)
				continue
			}

			//已经与客户端建立连接, 做点业务
			go func(){
				for {
					buf := make([]byte,512)
					cnt,err := conn.Read(buf)
					if err!=nil{
						fmt.Println("recv buf err ", err)
						continue
					}

					fmt.Printf("recv client buf %s,cnt %d\n", buf,cnt)

					// 回显
					if _,err := conn.Write(buf[:cnt]);err!=nil{
						fmt.Println("write back buf err",err)
						continue
					}

				}
			}()
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

func NewServer(name string)ziface.IServer{
	return &Server{
		Name:name,
		IPVersion: "tcp4",
		IP: "0.0.0.0",
		Port:9000,
	}
}