package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinxAll/zinx/utils"
	"zinxAll/zinx/ziface"
)

type DataPack struct {

}

func NewDataPack()*DataPack{
	return &DataPack{}
}


func (dp *DataPack) GetHeadLen() uint32 {
	return 8
}

//封包
// |datalen|msgID|data|
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 1.创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	// 2.将dataLen写入databuff中
	if err := binary.Write(dataBuff,binary.LittleEndian,msg.GetMsgLen());err!=nil{
		return nil,err
	}
	// 3.将MsgId 写进databuff中
	if err := binary.Write(dataBuff,binary.LittleEndian,msg.GetMsgId());err!=nil{
		return nil,err
	}
	// 4.将data数据 写进databuff中
	if err := binary.Write(dataBuff,binary.LittleEndian,msg.GetData());err!=nil{
		return nil,err
	}
	return dataBuff.Bytes(),nil
}

// 拆包方法(将包的Head信息都出来) 之后再根据head信息里的data的长度，再进行一次都
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewReader(binaryData)
	//只解压head信息，得到datalen和MsgID
	msg := &Message{}

	// 读dataLen
	if err :=binary.Read(dataBuff,binary.LittleEndian,&msg.DataLen);err!=nil{
		return nil,err
	}
	//  读MsgID
	if err:= binary.Read(dataBuff,binary.LittleEndian,&msg.Id);err!=nil{
		return nil,err
	}
	//判断datalen是否已经超出了我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize>0 && msg.DataLen>utils.GlobalObject.MaxPackageSize{
		return nil,errors.New(" too Large msg data recv! ")
	}

	return msg ,nil
}
