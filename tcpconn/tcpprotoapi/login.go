package tcpprotoapi

import (
	"bytes"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"net"
	"servicegunbattle/awesomeProject/bean"
	"servicegunbattle/awesomeProject/lib/logs"
	"servicegunbattle/awesomeProject/pb"
)

func Loginapi(cmd *pb.TcpLogin, cli *bean.Client,con net.Conn) {

	bean.ClientMaps.SetClient(con, *cli)
	bean.OnlineUIDs[cli.Uid] = 1
	//写入客户端消息
	sendmes := &pb.TcpResponseLogin{
		Result: proto.Bool(true),
		Uid: proto.Int32((rand.Int31n(10000)+10000)),
		UdpPort: proto.Int32(int32(10011)),
		Reconnect: proto.Bool(false),
	}
	senddata, _ := proto.Marshal(sendmes)
	logs.Info("Marshal ...",sendmes,"senddata长度:",len(senddata))
	newdata := make([]byte,0,len(senddata)+3)
	var inta int8 = 1
	bytesBuffer1 := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer1, binary.LittleEndian, &inta)

	var intb int16 = int16(len(senddata)+3)
	bytesBuffer2 := bytes.NewBuffer([]byte{})
	//binary.Write(bytesBuffer2, binary.BigEndian, &intb)
	binary.Write(bytesBuffer2, binary.LittleEndian, &intb)
	logs.Info("intb ...",intb)

	newdata = append(newdata,bytesBuffer1.Bytes()...)
	newdata = append(newdata,bytesBuffer2.Bytes()...)
	newdata = append(newdata,senddata...)

	cli.WriteClient(newdata)



}
