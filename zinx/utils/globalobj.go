package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinxAll/zinx/ziface"
)

type GlobalObj struct {
	TcpServer ziface.IServer  // zinx全局Server对象
	Host string
	TcpPort int  //当前服务器主机监听端口号
	Name string

	//
	Version string //当前zinx的版本号
	MaxConn int // 当前服务器主机允许的最大链接数
	MaxPackageSize uint32 // 数据包最大值

}

var (
	GlobalObject *GlobalObj
)

func (g *GlobalObj)Reload(){
	data ,err := ioutil.ReadFile("conf/zinx.json")
	if err!=nil{
		panic(any(err))
	}

	err = json.Unmarshal(data,&GlobalObject)
	if err!=nil{
		panic(any(err))
	}
}



func init(){
	GlobalObject = &GlobalObj{
		Name :"ZinxServerApp",
		Version:"V0.4",
		TcpPort: 9000,
		Host: "0.0.0.0",
		MaxConn: 1000,
		MaxPackageSize: 4096,
	}

	// 从conf/zinx.json中加载
	GlobalObject.Reload()

}


