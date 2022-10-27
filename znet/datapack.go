package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

type DataPack struct{}

func NewDP() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	//DataLen uint32(4字节)+ID uint32（4字节）
	return 8
}

func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	bytesbuf := bytes.NewBuffer([]byte{})

	//将data长度封装进去
	if err := binary.Write(bytesbuf, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	//将data id 封装进去
	if err := binary.Write(bytesbuf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//将data 消息本体封装进去
	if err := binary.Write(bytesbuf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return bytesbuf.Bytes(), nil
}

//拆包方法 （将包的head信息读取）之后再根据head里的len信息读取信息
func (d *DataPack) UnPack(data []byte) (ziface.IMessage, error) {
	//创建一个从二进制读取数据的ioReader
	databuf := bytes.NewReader(data)

	//只解压head信息的到len和id
	msg := &Message{}
	//读取datalen
	if err := binary.Read(databuf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读取id
	if err := binary.Read(databuf, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断datalen是否过长超出允许长度
	if utils.GlobalConfig.MaxPackageSize > 0 && msg.DataLen > utils.GlobalConfig.MaxPackageSize {
		return nil, errors.New("too large msg data recv!")
	}

	return msg, nil
}
