package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
	"servicegunbattle/awesomeProject/lib/logs"
	"servicegunbattle/awesomeProject/pb"
	"strings"
)

func main() {
	//open connection:
	conn, err := net.Dial("tcp", "192.168.3.13:1080")
	if err != nil {
		fmt.Println("Error dial:", err.Error())
		return
	}

	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("Please input your name:")
	//clientName, _ := inputReader.ReadString('\n')
	//inputClientName := strings.Trim(clientName,"\n")

	//send info to server until Quit
	_sendmes := pb.RequestCmd{
		Simple: &pb.ReqSimple{
			Tag: pb.ReqSimple_REQ_LOGIN_GAME.Enum(),
		},
		Reqlogin: &pb.TcpLogin{
			Userid:  proto.Int32(10001),
			Token: proto.String("token_10001"),
		},
	}
	senddata, err := proto.Marshal(&_sendmes)
	logs.Info("_sendmes =>",&_sendmes,"Marshal:",senddata)
	//int 类型转换成 若干长度的[]byte
	var sendLen int32 = int32(len(senddata))
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, &sendLen)
	newSend := make([]byte,0,len(senddata)+4)
	newSend = append(newSend,bytesBuffer.Bytes()...)
	newSend = append(newSend,senddata...)
	for {
		logs.Info("newSend:",newSend)

		fmt.Println("What do you send to the server? Type Q to quit.")
		content, _ := inputReader.ReadString('\n')
		inputContent := strings.Trim(content,"\n")
		if inputContent == "Q" {
			return
		}

		n, err := conn.Write(newSend)
		if err != nil {
			fmt.Println("Error Write:", err.Error())
			return
		}else {
			fmt.Println("Write-----n:", n)
		}
	}
}
