package ziface


type IConnManager interface {
	// 添加链接
	Add(conn IConnection)
	Remove(conn IConnection)
	Get(connID uint32)(IConnection,error)
	Len() int
	ClearConn()
}
