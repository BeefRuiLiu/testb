package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
	"runtime"
	"servicegunbattle/awesomeProject/base"
	"servicegunbattle/awesomeProject/game"
	"servicegunbattle/awesomeProject/game/room"
	"servicegunbattle/awesomeProject/lib/logs"
	"servicegunbattle/awesomeProject/lib/util"
	"servicegunbattle/awesomeProject/pb"
	"servicegunbattle/awesomeProject/tcpconn"
	"strconv"
	"time"
)
//用户信息
type User struct {
	userName string					//用户名
	userAddr *net.UDPAddr			//用户的监听连接地址
	userListenConn *net.UDPConn  	//用户的监听连接
	chatToConn *net.UDPConn 		//用户的聊天监听连接
}

//服务器监听端口
const LISTENPORT = 10011
//缓冲区
const BUFFSIZE = 1024
var allbuff = make([]byte, BUFFSIZE)
//在线用户
var onlineUser = make([]User, 0)
//在线状态判断缓冲区
var onlineCheckAddr = make([]*net.UDPAddr, 0)

//错误处理
func HandleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
}
func Cycle() {

	go func() {
		for {
			select {
			case <-time.After(time.Second * 60 * 1):
				fmt.Println("打印goroutine数量=========>>", runtime.NumGoroutine())
			}

		}
	}()
}
//消息处理
func oldHandleMessage(udpListener *net.UDPConn) {


	n, addr, err := udpListener.ReadFromUDP(allbuff)
	fmt.Println("服务数据长度 ：",n)
	fmt.Println("addr ：",addr)
	fmt.Println("服务数据==> ：",string(allbuff[:n]))
	HandleError(err)


	if n > 0 {
		//消息解析，[]byte -> []string
		msg := AnalyzeMessage(allbuff, n)

		switch msg[0] {
		//连接信息
		case "connect  ":
			//获取昵称+端口

			userName := msg[1]
			userListenPort := msg[2]
			//获取用户ip
			ip := AnalyzeMessage([]byte(addr.String()), len(addr.String()))
			//显示登录信息
			fmt.Println(" 昵称:", userName, " 地址:", ip[0], " 用户监听端口:", userListenPort, " 登录成功！")
			//创建对用户的连接，用于消息转发
			userAddr, err := net.ResolveUDPAddr("udp4", ip[0] + ":" + userListenPort)
			HandleError(err)

			userConn, err := net.DialUDP("udp4", nil, userAddr)
			HandleError(err)

			//因为连接要持续使用，不能在这里关闭连接
			//defer userConn.Close()
			//添加到在线用户
			onlineUser = append(onlineUser, User{userName, addr, userConn, nil})
			fmt.Println("用户加入在线数组：",onlineUser)

		case "online   ":
			//收到心跳包
			//fmt.Println("收到心跳包",addr)
			onlineCheckAddr = append(onlineCheckAddr, addr)
			//fmt.Println("收到心跳包",len(onlineCheckAddr))

		case "outline  ":
			//退出消息，未实现
			fmt.Println("收到退出消息",addr)
		case "chat     ":
			//会话请求


			//for i := 0; i < len(onlineUser); i++ {
			//	fmt.Println("onlineUser[i]:",onlineUser[i].userName)
			//	//if onlineUser[i].userName == msg[1] {
			//	//	index = i
			//	//}
			//	onlineUser[i].userListenConn.Write([]byte(onlineUser[i].userName + " " + time.Now().String() +"\n" + msg[1]))
			//
			//}

			for _, v := range onlineUser {
					//writerdata := []byte(v.userName + " " + time.Now().String() +"\n" + msg[1])
					//writerdata = append(writerdata,data...)
					//v.userListenConn.Write(writerdata)
				ack := pb.CommandCode_ACK
				_c := &pb.CommandHeader{
					CommandCode:  &ack,
					Flag: proto.Int32(2),
				}
				data, _ := proto.Marshal(_c)
					v.userListenConn.Write(data)
			}

			//寻找请求对象
			//index := -1
			//将所请求对象的连接添加到请求者中
			//if index != -1 {
			//	nowUser, _ := FindUser(addr)
			//	onlineUser[nowUser].chatToConn = onlineUser[index].userListenConn
			//}
		case "get      ":
			//向请求者返回在线用户信息
			index, _ := FindUser(addr)
			if index < 0 {
				fmt.Println("FindUser(addr)未找到该用户")
				return
			}
			onlineUser[index].userListenConn.Write([]byte("当前共有" + strconv.Itoa(len(onlineUser)) + "位用户在线"))
			for i, v := range onlineUser {
				onlineUser[index].userListenConn.Write([]byte("" + strconv.Itoa(i + 1) + ":" + v.userName))
			}
		default:
			//消息转发
			//获取当前用户
			sendIndex, _ := FindUser(addr)
			nowTime := time.Now()
			nowHour := strconv.Itoa(nowTime.Hour())
			nowMinute := strconv.Itoa(nowTime.Minute())
			nowSecond := strconv.Itoa(nowTime.Second())
			for index := 0 ; index < len(onlineUser); index ++ {
				//获取时间
				//请求会话对象是否存在
				if onlineUser[index].userListenConn == nil {
					onlineUser[index].userListenConn.Write([]byte("对方不在线"))
				} else {
					onlineUser[index].userListenConn.Write([]byte(onlineUser[sendIndex].userName + " " + nowHour + ":" + nowMinute + ":" + nowSecond + "\n" + msg[0]))
				}
			}


		}
	}
}
func HandleMessage(n int,addr *net.UDPAddr,buff []byte) {


	//n, addr, err := udpListener.ReadFromUDP(buff)
	//fmt.Println("服务数据长度 ：",n)
	//fmt.Println("addr ：",addr)
	//fmt.Println("服务数据==> ：",string(buff[:n]))
	//HandleError(err)


	if n > 0 {
		//消息解析，[]byte -> []string
		msg := AnalyzeMessage(buff, n)
		switch msg[0] {
		//连接信息
		case "connect  ":
			//获取昵称+端口

			userName := msg[1]
			userListenPort := msg[2]
			//获取用户ip
			ip := AnalyzeMessage([]byte(addr.String()), len(addr.String()))
			//显示登录信息
			fmt.Println(" 昵称:", userName, " 地址:", ip[0], " 用户监听端口:", userListenPort, " 登录成功！")
			//创建对用户的连接，用于消息转发
			userAddr, err := net.ResolveUDPAddr("udp4", ip[0] + ":" + userListenPort)
			HandleError(err)

			userConn, err := net.DialUDP("udp4", nil, userAddr)
			HandleError(err)

			//因为连接要持续使用，不能在这里关闭连接
			//defer userConn.Close()
			//添加到在线用户
			onlineUser = append(onlineUser, User{userName, addr, userConn, nil})
			fmt.Println("用户加入在线数组：",onlineUser)

		case "online   ":
			//收到心跳包
			//fmt.Println("收到心跳包",addr)
			onlineCheckAddr = append(onlineCheckAddr, addr)
			//fmt.Println("收到心跳包",len(onlineCheckAddr))

		case "outline  ":
			//退出消息，未实现
			fmt.Println("收到退出消息",addr)
		case "chat     ":
			//会话请求


			//for i := 0; i < len(onlineUser); i++ {
			//	fmt.Println("onlineUser[i]:",onlineUser[i].userName)
			//	//if onlineUser[i].userName == msg[1] {
			//	//	index = i
			//	//}
			//	onlineUser[i].userListenConn.Write([]byte(onlineUser[i].userName + " " + time.Now().String() +"\n" + msg[1]))
			//
			//}

			for _, v := range onlineUser {
					//writerdata := []byte(v.userName + " " + time.Now().String() +"\n" + msg[1])
					//writerdata = append(writerdata,data...)
					//v.userListenConn.Write(writerdata)
				ack := pb.CommandCode_ACK
				_c := &pb.CommandHeader{
					CommandCode:  &ack,
					Flag: proto.Int32(2),
				}
				data, _ := proto.Marshal(_c)
					v.userListenConn.Write(data)
			}

			//寻找请求对象
			//index := -1
			//将所请求对象的连接添加到请求者中
			//if index != -1 {
			//	nowUser, _ := FindUser(addr)
			//	onlineUser[nowUser].chatToConn = onlineUser[index].userListenConn
			//}
		case "get      ":
			//向请求者返回在线用户信息
			index, _ := FindUser(addr)
			if index < 0 {
				fmt.Println("FindUser(addr)未找到该用户")
				return
			}
			onlineUser[index].userListenConn.Write([]byte("当前共有" + strconv.Itoa(len(onlineUser)) + "位用户在线"))
			for i, v := range onlineUser {
				onlineUser[index].userListenConn.Write([]byte("" + strconv.Itoa(i + 1) + ":" + v.userName))
			}
		default:
			//消息转发
			//获取当前用户
			sendIndex, _ := FindUser(addr)
			nowTime := time.Now()
			nowHour := strconv.Itoa(nowTime.Hour())
			nowMinute := strconv.Itoa(nowTime.Minute())
			nowSecond := strconv.Itoa(nowTime.Second())
			for index := 0 ; index < len(onlineUser); index ++ {
				//获取时间
				//请求会话对象是否存在
				if onlineUser[index].userListenConn == nil {
					onlineUser[index].userListenConn.Write([]byte("对方不在线"))
				} else {
					onlineUser[index].userListenConn.Write([]byte(onlineUser[sendIndex].userName + " " + nowHour + ":" + nowMinute + ":" + nowSecond + "\n" + msg[0]))
				}
			}


		}
	}
}
//消息处理
func NewHandleMessage(udpListener *net.UDPConn) {


	n, addr, err := udpListener.ReadFromUDP(allbuff)
	fmt.Println("服务数据长度 ：",n)
	fmt.Println("addr ：",addr)
	fmt.Println("服务数据==> ：",string(allbuff[:n]))
	HandleError(err)

	if n > 0 {
		_CommandHeader := &pb.CommandHeader{
		}
		err = proto.Unmarshal(allbuff[:n],_CommandHeader)

		//fmt.Println("test Unmarshal(data) ...",_CommandHeader)
		if err != nil {
			switch _CommandHeader.GetCommandCode() {
			//连接信息
			case pb.CommandCode_ACK:
				//获取昵称+端口

				userName := _CommandHeader.String()
				userListenPort := _CommandHeader.String()
				//获取用户ip
				ip := AnalyzeMessage([]byte(addr.String()), len(addr.String()))
				//显示登录信息
				fmt.Println(" 昵称:", userName, " 地址:", ip[0], " 用户监听端口:", userListenPort, " 登录成功！")
				//创建对用户的连接，用于消息转发
				userAddr, err := net.ResolveUDPAddr("udp4", ip[0] + ":" + userListenPort)
				HandleError(err)

				userConn, err := net.DialUDP("udp4", nil, userAddr)
				HandleError(err)

				//因为连接要持续使用，不能在这里关闭连接
				//defer userConn.Close()
				//添加到在线用户
				onlineUser = append(onlineUser, User{userName, addr, userConn, nil})
				fmt.Println("用户加入在线数组：",onlineUser)

			case pb.CommandCode_ACKCONNECT:
				//收到心跳包
				//fmt.Println("收到心跳包",addr)
				onlineCheckAddr = append(onlineCheckAddr, addr)
				//fmt.Println("收到心跳包",len(onlineCheckAddr))
			case pb.CommandCode_CONNECT:

				for _, v := range onlineUser {
					//writerdata := []byte(v.userName + " " + time.Now().String() +"\n" + msg[1])
					//writerdata = append(writerdata,data...)
					//v.userListenConn.Write(writerdata)
					ack := pb.CommandCode_ACK
					_c := &pb.CommandHeader{
						CommandCode:  &ack,
						Flag: proto.Int32(2),
					}
					data, _ := proto.Marshal(_c)
					v.userListenConn.Write(data)
				}

			default:
				//消息转发
				//获取当前用户
				sendIndex, _ := FindUser(addr)
				nowTime := time.Now()
				nowHour := strconv.Itoa(nowTime.Hour())
				nowMinute := strconv.Itoa(nowTime.Minute())
				nowSecond := strconv.Itoa(nowTime.Second())
				for index := 0 ; index < len(onlineUser); index ++ {
					//获取时间
					//请求会话对象是否存在
					if onlineUser[index].userListenConn == nil {
						onlineUser[index].userListenConn.Write([]byte("对方不在线"))
					} else {
						onlineUser[index].userListenConn.Write([]byte(onlineUser[sendIndex].userName + " " + nowHour + ":" + nowMinute + ":" + nowSecond + "\n"))
					}
				}


			}
		}else {
			fmt.Println("test err ...",err)
		}

	}
}
//消息解析，[]byte -> []string
func AnalyzeMessage(buff []byte, len int) ([]string) {
	//fmt.Println("消息解析buff==>>",string(buff[:len]))
	analMsg := make([]string, 0)
	strNow := ""
	for i := 0; i < len; i++ {
		if string(buff[i:i + 1]) == ":" {
			analMsg = append(analMsg, strNow)
			strNow = ""
		} else {
			strNow += string(buff[i:i + 1])
		}
	}
	analMsg = append(analMsg, strNow)
	return analMsg
}
//寻找用户，返回（位置，是否存在）
func FindUser(addr *net.UDPAddr) (int, bool) {
	alreadyhave := false
	index := -1
	for i := 0; i < len(onlineUser); i++ {

		if onlineUser[i].userAddr.String() == addr.String() {
			alreadyhave = true
			index = i
			break
		}
	}
	return index, alreadyhave
}
//处理用户在线信息（暂时仅作删除用户使用）
func HandleOnlineMessage(addr *net.UDPAddr, state bool) {
	index, alreadyhave := FindUser(addr)
	if state == false {
		if alreadyhave {
			onlineUser = append(onlineUser[:index], onlineUser[index + 1:len(onlineUser)]...)
		}
	}
}
//在线判断，心跳包处理，每5s查看一次所有已在线用户状态
func OnlineCheck() {
	for {
		onlineCheckAddr = make([]*net.UDPAddr, 0)
		sleepTimer := time.NewTimer(time.Second * 5)
		<- sleepTimer.C
		//fmt.Println(time.Now().String(),"在线判断，心跳包处理，每5s查看一次所有已在线用户状态")
		for i := 0; i < len(onlineUser); i++ {
			haved := false

		FORIN:for j := 0; j < len(onlineCheckAddr); j++ {
			if onlineUser[i].userAddr.String() == onlineCheckAddr[j].String() {
				haved = true
				break FORIN
			}
		}
			if !haved {
				fmt.Println(onlineUser[i].userAddr.String() + "退出！")
				HandleOnlineMessage(onlineUser[i].userAddr, false)
				i--
			}

		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	base.InitReadKeyMapData()
	//监听TCP端口
	go tcpconn.Server()
	game.InitGameMes()
	//监听地址
	udpAddr, err := net.ResolveUDPAddr("udp4", "192.168.3.13:"+strconv.Itoa(LISTENPORT))
	HandleError(err)
	//监听连接
	udpListener, err := net.ListenUDP("udp4", udpAddr)
	logs.Info("=====================UDP开始监听：%v", udpAddr.String())


	HandleError(err)

	defer udpListener.Close()

	//在线状态判断
	go OnlineCheck()
	//go Cycle()
	buf := make([]byte, 1024)
	//raddr := net.UDPAddr{
	//	IP:   net.IPv4(255, 255, 255, 255),
	//	//Port: LISTENPORT + 1,
	//
	//}

	for {
		n, ctlAddr, err := udpListener.ReadFromUDP(buf)
		if err != nil {
			logs.Info("%v", err)
			continue
		}

		logs.Info("收到UDP消息---ctlAddr:", ctlAddr.String())

		csid,err := util.BytesToIntU(buf[:1])
		if err != nil {
			logs.Info("BytesToIntU--err:",err.Error())
		}
		//scidstr := strconv.Itoa(csid)

		logs.Info("读取到的消息ID：%v",csid)
		//TODO--需要解析proto
		senddata := make([]byte,0)
		UDPAddrs := make([]*net.UDPAddr,0)
		if csid == 51 {
			_udpcmd := &pb.UdpBattleReady{
			}
			err = proto.Unmarshal(buf[3:n], _udpcmd)
			//_, err = udpListener.Write([]byte("广播消息:" + ctlAddr.String() + string(buf[:n])))
			//_, err = udpListener.Write(buf[:n])
			isbegin := false
			//UDPAddrs := make([]*net.UDPAddr,0)
			a :=make([]int,0)
			b :=make([]int,0)
			game.GameMes.RoomsMap.ReadRange(func(roomid int, _room *room.Room) {
				logs.Info("readRange--AllUIDs",_room.AllUIDs)
				if util.HasElement(_room.AllUIDs,int(_udpcmd.GetUid())) {
					_room.Playermap.ForEach(func(key int, s *room.Player) {
						if key == int(_udpcmd.GetUid()) {
							s.UDPaddress = ctlAddr
						}
						UDPAddrs = append(UDPAddrs,s.UDPaddress)
					})
					if !util.HasElement(_room.ReadyUIDs,int(_udpcmd.GetUid())) {
						_room.ReadyUIDs = append(_room.ReadyUIDs,int(_udpcmd.GetUid()))
						//return
					}
					a = _room.ReadyUIDs
					b = _room.AllUIDs

					if len(_room.ReadyUIDs) == len(_room.AllUIDs) && len(_room.AllUIDs) > 1 {
						isbegin = true

					}
				}

			})
			if !isbegin {
				logs.Info("等待房间用户准备room.ReadyUIDs",a)
				logs.Info("等待房间用户准备room.alluids",b)
				continue
			}
			sendmes := &pb.UdpBattleStart{
			}
			logs.Info("UDPAddrs-->:%v",UDPAddrs)
			senddata, _ = proto.Marshal(sendmes)
		}else if csid == 53 {
			_udpcmd := &pb.UdpUpPlayerOperations{
			}
			FrameID := 1
			err = proto.Unmarshal(buf[3:n], _udpcmd)
			logs.Info("读取_udpcmd:",_udpcmd.GetOperation())
			operation := make([]*pb.PlayerOperation,0)

			isbegin := false
			game.GameMes.RoomsMap.ReadRange(func(roomid int, _room *room.Room) {

				_room.Playermap.ForEach(func(uid int, s *room.Player) {

					logs.Info("s.UDPaddress:",s.UDPaddress.String())
					logs.Info("ctlAddr:",ctlAddr.String())
					if s.UDPaddress.String() == ctlAddr.String() {
						s.PlayerOpt = _udpcmd.GetOperation()
					}

					UDPAddrs = append(UDPAddrs,s.UDPaddress)
					if s.PlayerOpt.GetBattleID() > 0 {
						operation = append(operation,s.PlayerOpt)
					}
				})

				logs.Info("operation:",operation)
				if len(operation) == len(_room.AllUIDs) {
					FrameID = _room.FrameID
					isbegin = true
					_room.FrameID += 1
					_room.Playermap.ForEach(func(uid int, s *room.Player) {
						s.PlayerOpt = &pb.PlayerOperation{}
					})
				}
			})
			if !isbegin {
				logs.Info("等待房间其他用户操作",operation)
				continue
			}

			logs.Info("完成用户操作",operation)
			//operation = append(operation,_udpcmd.GetOperation())
			operations := &pb.AllPlayerOperation{
				Operations: operation,
			}
			sendmes := &pb.UdpDownFrameOperations{
				FrameID: proto.Int32(int32(FrameID)),
				Operations: operations,
			}
			//FrameID ++
			logs.Info("sendmes:%v",sendmes)
			senddata, _ = proto.Marshal(sendmes)
		}else if csid == 55 {
			_udpcmd := &pb.UdpUpDeltaFrames{
			}
			err = proto.Unmarshal(buf[3:n], _udpcmd)
			//_, err = udpListener.Write([
			//]byte("广播消息:" + ctlAddr.String() + string(buf[:n])))
			//_, err = udpListener.Write(buf[:n])

			sendmes := &pb.UdpDownDeltaFrames{

			}
			logs.Info("proto.UdpUpDeltaFrames-->:%v",_udpcmd)
			senddata, _ = proto.Marshal(sendmes)
		}else if csid == 57 {
			_udpcmd := &pb.UdpUpGameOver{
			}
			err = proto.Unmarshal(buf[3:n], _udpcmd)
			//_, err = udpListener.Write([
			//]byte("广播消息:" + ctlAddr.String() + string(buf[:n])))
			//_, err = udpListener.Write(buf[:n])
			removeRoomID := 0
			game.GameMes.RoomsMap.ReadRange(func(roomid int, _room *room.Room) {
					_room.Playermap.ForEach(func(key int, s *room.Player) {
							removeRoomID = roomid
						UDPAddrs = append(UDPAddrs,s.UDPaddress)
					})


			})
			game.GameMes.RoomsMap.RemoveRoom(removeRoomID)
			sendmes := &pb.UdpDownGameOver{

			}
			logs.Info("proto.UdpUpGameOver-->:%v",_udpcmd)
			senddata, _ = proto.Marshal(sendmes)
		}else {
			logs.Info("csid不存在-->:%v",csid)
		}
		//_udpcmd := &pb.UdpBattleReady{
		//}
		//err = proto.Unmarshal(buf[3:n], _udpcmd)

		//logs.Info("proto.Unmarshal-->:%v",_udpcmd)

		if true {
			newdata := make([]byte,0,len(senddata)+3)
			var inta int8 = 51
			if csid == 53 {
				inta = 53
			}else if csid == 55 {
				inta = 55
			}else if csid == 57 {
				inta = 57
			}
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
			//_sendtomes := &pb.UdpHeader{
			//
			//}
			//sendtodata, _ := proto.Marshal(_sendtomes)

			//_, err = udpListener.WriteTo([]byte("server转发广播消息:"+ctlAddr.String()+string(buf[:n])), &raddr)
			for _, v := range UDPAddrs {
				_, err = udpListener.WriteTo(newdata, v)
				if err != nil {
					logs.Info("发送UDP消息出错:err%v", err)
				}else {
					logs.Info("发送UDP消息成功",newdata,"To-》",v)
				}
			}

			//消息处理

			////HandleMessage(udpListener)
		}else {
			logs.Info("senddata==>nil")
		}

	}
}