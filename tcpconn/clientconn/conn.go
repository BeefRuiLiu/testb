package clientconn

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"net"
	"runtime"
	"servicegunbattle/awesomeProject/bean"
	"servicegunbattle/awesomeProject/game"
	"servicegunbattle/awesomeProject/game/room"
	"servicegunbattle/awesomeProject/lib/logs"
	"servicegunbattle/awesomeProject/lib/util"
	"servicegunbattle/awesomeProject/pb"
	"servicegunbattle/awesomeProject/tcpconn/tcpprotoapi"
	"time"
)


func Cycle() {

	//go func() {
	//	for {
	//		select {
	//		case <-time.After(time.Second * 10):
	//			//delkey := make([]net.Conn,0)
	//			//logs.Info( "心跳检测===============>>")
	//
	//			bean.ClientMaps.Range(func(key net.Conn, client bean.Client) {
	//				//log.Println("遍历map==>key:",key,"val:",client)
	//				tag := time.Now().Unix() - client.Heartbeat
	//				if tag > 90 {
	//					logs.Info(  "90秒未收到客户端的心跳包断开连接===============>>%v", client.Uid)
	//					//delkey = append(delkey,key)
	//					closeMes := new(bean.ReadMessage)
	//					closeMes.Code = -1
	//					closeMes.Currencyid = client.Uid
	//					bean.Broadcast <- closeMes
	//
	//				}
	//			})
	//		}
	//	}
	//}()
	go func() {
		for {
			select {
			case <-time.After(time.Second * 60 * 5):
				logs.Info( "在线用户数量===============>>%v", len(bean.ClientMaps.ClientsMap))
				logs.Info(  "打印goroutine数量=========>>%v", runtime.NumGoroutine())
			}
		}
	}()
	go func() {
		for {
			select {
			case <-time.After(time.Second * 15 ):
				game.GameMes.RoomsMap.ReadRange(func(roomid int, room *room.Room) {
					logs.Info(  "roomid=========>>",roomid,"room.Userinfos数量=========>>",len(room.Userinfos))
				})

			}
		}
	}()
}

//全局消息
//func HandleMessages() {
//	for {
//		select {
//		case msg, ok := <-bean.Broadcast:
//			if !ok {
//				logs.Info( "-----------bean.Broadcast关闭")
//				return
//			} else {
//				issend := 0
//				closeConn := make([]net.Conn, 0)
//				bean.ClientMaps.Range(func(conn net.Conn, client bean.Client) {
//					//msg.RemoteAddr != conn.RemoteAddr().String() ||
//					if true {
//						//把当前的信息传递给通讯C,此时map中的c的数据相同
//						if msg.Code == baseim.Code_RoomChat {
//							if msg.Currencyid == client.Roomid {
//								client.C <- msg
//								issend = 1
//								//return false
//							}
//						} else if msg.Code == baseim.Code_FriendChat {
//							if msg.Currencyid == client.Uid {
//								client.C <- msg
//								issend = 1
//								return
//							}
//						} else if msg.Code == baseim.Code_Login {
//							//log.Println("======mes",msg)
//							if msg.Uid == client.Uid {
//								client.C <- msg
//								issend = 1
//								return
//							}
//						} else if msg.Code == baseim.Code_Heartbeat {
//							//发送心跳
//							if msg.Uid == client.Uid {
//								client.C <- msg
//								issend = 1
//								return
//							}
//
//						} else {
//							logs.Info(  "用户：", msg.Uid, "消息Code错误=======>", msg.Code)
//							return
//						}
//
//					}
//				})
//				for _, v := range closeConn {
//					bean.ClientMaps.DelClient(v)
//				}
//				if issend == 1 {
//					log.Printf("打印broadcast已经发送的msg: %v", msg)
//				}
//				//TODO----------保存到数据库
//				//if msg.Sockettype == 1 && (msg.Code == baseim.Code_RoomChat || msg.Code == baseim.Code_FriendChat) && msg.Data != "" {
//				//	//加密保存
//				//	strEncrypted, deserr := tools.Encrypt(msg.Data, tools.DESKey)
//				//	if deserr != nil {
//				//		base.Err(ImTag,"加密失败:",deserr.Error())
//				//	}else {
//				//		msg.Issend = issend
//				//		msg.Data = strEncrypted
//				//		controllers.Model.InsertUserMes(msg)
//				//	}
//				//
//				//}
//			}
//		}
//	}
//}



//给客户端发送消息
func writer(cli bean.Client, conn net.Conn) {
	defer conn.Close()
	for {
		select {
		case res, ok := <-cli.C:
			if !ok {
				delete(bean.OnlineUIDs,cli.Uid)
				logs.Info( "-----------客户端消息通道关闭---结束goroutine")
				return
			} else {
				//senddata, err := proto.Marshal(res)
				//if err != nil {
				//	//通道关闭通知-1
				//	delete(bean.OnlineUIDs,cli.Uid)
				//	close(cli.C)
				//	logs.Info(  "-----------客户端消息通道关闭消息---序列化写入proto消息错误",err.Error())
				//	return
				//}
				////int 类型转换成 若干长度的[]byte
				//var sendLen int32 = int32(len(senddata))
				//bytesBuffer := bytes.NewBuffer([]byte{})
				//binary.Write(bytesBuffer, binary.BigEndian, &sendLen)
				//newSend := make([]byte,0,len(senddata)+4)
				//newSend = append(newSend,bytesBuffer.Bytes()...)
				//newSend = append(newSend,senddata...)
				//发送消息
				_, err := conn.Write(res)
				if err != nil {
					logs.Info( "给客户端发送消息失败:", err.Error())
				} else {
					logs.Info( "给客户端发送消息成功")
				}
			}
		}

	}

}

//断开连接
func disClientConn(conn net.Conn) {
	//TODO-------------断开连接记录玩家下线时间
	//RemoteAddr := conn.RemoteAddr().String()
	//close(bean.ClientMaps.GetClient(conn).C)
	//conn.Close()
	bean.ClientMaps.DelClient(conn)

}

//读取客户端的消息
func HandleConnections(con net.Conn) {
	defer func() {
		if e := recover(); e != nil {
			trace := make([]byte, 1<<16)
			n := runtime.Stack(trace, true)
			logs.Info(  "%v", fmt.Errorf("panic recover\n %v\n stack trace %d bytes\n %s",
				e, n, trace[:n]))
		}
	}()
	logs.Info( "handleConnections--->客户端连接: ", con.RemoteAddr())

	//Get client's name

	defer con.Close()
	//defer log.Println("Client退出连接--------defer()HandleConnections:", con.RemoteAddr())

	//_reader := bufio.NewReaderSize(con, 2048)
	//go func() {
	//	select {
	//	case <-time.After(time.Second * 30):
	//		//con.Close()
	//		c := bean.ClientMaps.GetClient(con)
	//		if c.Uid == 0 {
	//			//	11s未发送登陆消息则断开连接
	//			logs.Info( "30s没有发送消息则断开连接 %v .", con.RemoteAddr())
	//			con.Close()
	//			return
	//		}
	//	}
	//}()
	cli := bean.NewClient(con.RemoteAddr().String())
	//写入玩家消息
	go writer(*cli, con)

	// Begin recieve message from client
	for {
		_reader := bufio.NewReaderSize(con, 10240)
		//data, err := codec.Decode(_reader)
		data := make([]byte, 2048)
		n, err := _reader.Read(data)
		logs.Info("_reader.Read(pack):",n)
		//if err != nil {
		//	logs.Info("读取消息的内容err1",err.Error())
		//	return
		//}
		if err == io.EOF {
			log.Println("Client （for）主动关闭通道:",err.Error())
			closeMes := new(bean.ReadMessage)
			closeMes.Code = -1
			closeMes.Currencyid = cli.Uid
			bean.Broadcast <- closeMes
			//cli.C <- closeMes
			//disClientConn(con)
			return
		}
		if data == nil {
			logs.Info( "Client 未读取到数据:")
			closeMes := new(bean.ReadMessage)
			closeMes.Code = -1
			closeMes.Currencyid = cli.Uid
			bean.Broadcast <- closeMes
			return
		}
		csid,err := util.BytesToIntU(data[:1])
		if err != nil {
			logs.Info("BytesToIntU--err:",err.Error())
		}
		//scidstr := strconv.Itoa(csid)

		logs.Info("读取到的消息ID：%v",csid)
		//logs.Info("scidstr",scidstr)
		if csid == int(pb.CSID_TCP_LOGIN.Number()) {
			tcpMes := &pb.TcpLogin{
			}
			potoerr := proto.Unmarshal(data[3:n],tcpMes)
			if potoerr != nil {
				log.Println("解析proto错误:",potoerr.Error())
				//log.Println("json解析错误:err",jsonerr)
				return
			}
			logs.Info("tcpMes===>",tcpMes)
			tcpprotoapi.Loginapi(tcpMes,cli,con)


		}else if csid == int(pb.CSID_TCP_REQUEST_MATCH.Number()) {
			tcpMes := &pb.TcpRequestMatch{
			}
			potoerr := proto.Unmarshal(data[3:n],tcpMes)
			if potoerr != nil {
				log.Println("解析proto错误:",potoerr.Error())
				//log.Println("json解析错误:err",jsonerr)
				return
			}
			logs.Info("tcpMes===>",tcpMes)

			//bean.ClientMaps.SetClient(con, *cli)
			//bean.OnlineUIDs[cli.Uid] = 1
			//写入客户端消息
			sendmes := &pb.TcpResponseRequestMatch{

			}
			senddata, _ := proto.Marshal(sendmes)
			logs.Info("Marshal ...",sendmes,"senddata长度:",len(senddata))
			newdata := make([]byte,0,len(senddata)+3)
			var inta int8 = 10
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
			logs.Info("Enter_ROOM-tcpMes-->",tcpMes)
			//
			tcpprotoapi.RequestMatch(tcpMes,cli,con)
		}else if csid == int(pb.CSID_TCP_CANCEL_MATCH.Number()) {
			tcpMes := &pb.TcpCancelMatch{
			}
			potoerr := proto.Unmarshal(data[3:n],tcpMes)
			if potoerr != nil {
				log.Println("解析proto错误:",potoerr.Error())
				//log.Println("json解析错误:err",jsonerr)
				return
			}
			logs.Info("tcpMes==>>",tcpMes)
			//bean.ClientMaps.SetClient(con, *cli)
			//bean.OnlineUIDs[cli.Uid] = 1
			//写入客户端消息
			sendmes := &pb.TcpResponseCancelMatch{

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
			logs.Info("Enter_ROOM-tcpMes-->",tcpMes)

		}else {
			logs.Info("消息ID不存在======!=",csid)
			closeMes := new(bean.ReadMessage)
			closeMes.Code = -1
			closeMes.Currencyid = cli.Uid
			bean.Broadcast <- closeMes
			//cli.C <- closeMes
			//disClientConn(con)
			return
		}
		//[]byte ==》proto消息结构
		//tcpMes := &pb.RequestCmd{
		//}
		//potoerr := proto.Unmarshal(data,tcpMes)

		//if potoerr != nil {
		//	log.Println("解析proto错误:",potoerr.Error())
		//	//log.Println("json解析错误:err",jsonerr)
		//}else {
		//	//TODO ---- 根据简单消息判断
		//	switch tcpMes.GetSimple().GetTag() {
		//		case pb.ReqSimple_REQ_LOGIN_GAME:
		//			logs.Info("LOGIN_GAME-tcpMes-->",tcpMes)
		//			tcpprotoapi.Loginapi(tcpMes,cli,con)
		//		break
		//		case pb.ReqSimple_REQ_Enter_ROOM:
		//			logs.Info("Enter_ROOM-tcpMes-->",tcpMes)
		//			tcpprotoapi.EnterRoom(tcpMes,cli,con)
		//		break
		//		case pb.ReqSimple_REQ_Out_ROOM:
		//			logs.Info("Out_ROOM-tcpMes-->",tcpMes)
		//			tcpprotoapi.OutRoom(tcpMes,cli,con)
		//		break
		//		case pb.ReqSimple_REQ_Enter_Game:
		//			logs.Info("Enter_Game-tcpMes-->",tcpMes)
		//			tcpprotoapi.EnterGame(tcpMes,cli,con)
		//		break
		//	}
		//}
	}
}
