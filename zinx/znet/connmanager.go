package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinxAll/zinx/ziface"
)

type ConnManager struct{
	connections map[uint32] ziface.IConnection //管理的链接集合
	connLock sync.RWMutex
}

func NewConnManager() *ConnManager{
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源map,加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connID=" , conn.GetConnID(),"add to ConnManager successfully:conn num=", connMgr.Len())

}

func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	// 删除连接信息
	delete(connMgr.connections,conn.GetConnID())
}

func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn,ok := connMgr.connections[connID];ok{
		return conn,nil
	}else{
		return nil,errors.New("connection not FOUND! ")
	}
}

func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除conn并停止conn的工作
	for connID,conn := range connMgr.connections{
		conn.Stop()
		delete(connMgr.connections,connID)
	}

	fmt.Println("Clear All connections succ! conn num= ", connMgr.Len())

}

