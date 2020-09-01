package codec

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"servicegunbattle/awesomeProject/lib/logs"
)

func Encode(message string) ([]byte, error) {

	fmt.Println("Encode......")
	// 读取消息的长度
	var length int32 = int32(len(message))
	var pkg *bytes.Buffer = new(bytes.Buffer)
	// 写入消息头
	err := binary.Write(pkg, binary.LittleEndian, length)
	//fmt.Println("Encode...... length", length)
	//fmt.Println("Encode...... pkg.len", pkg.Len())
	if err != nil {
		fmt.Println("Encode......err1", err.Error())
		return nil, err
	}
	// 写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, []byte(message))
	if err != nil {
		fmt.Println("Encode......err2", err.Error())
		return nil, err
	}
	//fmt.Println("Encode...... pkg.Bytes()1", pkg.String())
	//fmt.Println("Encode...... pkg.Bytes()2", len(pkg.Bytes()))

	return pkg.Bytes(), nil
}

func Decode(reader *bufio.Reader) ([]byte, error) {
	// 读取消息的长度
	lengthByte, _ := reader.Peek(4)
	lengthBuff := bytes.NewBuffer(lengthByte)
	logs.Info("lengthByte",lengthByte)

	var length int32
	err := binary.Read(lengthBuff, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}
	logs.Info("length",length)

	if int32(reader.Buffered()) < length+4 {
		logs.Info("读取消息的内容err1",reader.Buffered(),"length+4",length+4)
		return nil, err
	}

	// 读取消息真正的内容
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		logs.Info("读取消息的内容err1")
		return nil, err
	}
	return pack[4:], nil
}
