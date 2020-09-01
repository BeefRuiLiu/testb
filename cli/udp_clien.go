package main

import (
	"bufio"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
	"servicegunbattle/awesomeProject/lib/logs"
	"servicegunbattle/awesomeProject/pb"
	"strconv"
	"time"
)
//数据包头，标识数据内容
var reflectString = map[string]string {
	"连接": 		"connect  :",
	"在线": 		"online   :",
	"聊天": 		"chat     :",
	"在线用户": 	"get      :",
}

//服务器端口
const CLIENTPORT = 10011
//数据缓冲区
const BUFFSIZE = 1024
var buff = make([]byte, BUFFSIZE)

//错误消息处理
func HandleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
}
//发送消息
func SendMessage(udpConn *net.UDPConn) {
	scaner := bufio.NewScanner(os.Stdin)

	for scaner.Scan() {
		if scaner.Text() == "exit" {
			return
		}
		fmt.Println("scaner.Text()",scaner.Text())
		udpConn.Write([]byte(reflectString["聊天"]+scaner.Text()))

	}
}
//接收消息
func HandleMessage(udpListener *net.UDPConn) {
	for {
		n, _, err := udpListener.ReadFromUDP(buff)
		HandleError(err)

		if n > 0 {
			_c2 := &pb.ResponseCmd{
			}
			err = proto.Unmarshal(buff[:n],_c2)
			if err != nil {
				fmt.Println("解析proto数据err ...",err)
			}

			fmt.Println("test Unmarshal(data) ...",_c2)

			fmt.Println("接收到服务端返回的消息:",_c2)
		}
	}
}
/*
func AnalyzeMessage(buff []byte, len int) ([]string) {
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
}*/
//发送心跳包
func SendOnlineMessage(udpConn *net.UDPConn) {
	for {
		//每间隔1s向服务器发送一次在线信息
		udpConn.Write([]byte(reflectString["在线"]))
		sleepTimer := time.NewTimer(time.Second*2)
		<- sleepTimer.C
	}
}

func main() {
	//判断命令行参数，参数应该为服务器ip
	addrss()
	//if len(os.Args) != 2 {
	//	fmt.Println("程序命令行参数错误！","len(os.Args)",len(os.Args),os.Args)
	//	os.Exit(2)
	//}
	//获取ip
	//host := os.Args[1]
	//host := base.SERVER_IPas
	//host := "127.0.0.1"
	host := "192.168.3.13"

	//udp地址
	udpAddr, err := net.ResolveUDPAddr("udp4", host + ":" + strconv.Itoa(CLIENTPORT))
	fmt.Println("udpAddr==",udpAddr.String())
	HandleError(err)
	//udp连接
	//laddr := net.UDPAddr{
	//	IP:   net.IPv4(192, 168, 3, 22),
	//	Port: 10006,
	//}
	udpConn, err := net.DialUDP("udp4", nil, udpAddr)
	HandleError(err)

	//本地监听端口
	//newSeed := rand.NewSource(int64(time.Now().Second()))
	//newRand := rand.New(newSeed)
	randPort := CLIENTPORT+1
	//randPort := newRand.Intn(30000) + 10007
	//本地监听udp地址
	//udpLocalAddr, err := net.ResolveUDPAddr("udp4", "192.168.3.13:" + strconv.Itoa(randPort))
	udpAdrss := addrss() +":"
	udpLocalAddr, err := net.ResolveUDPAddr("udp4", udpAdrss + strconv.Itoa(randPort))
	HandleError(err)

	//本地监听udp连接
	udpListener, err := net.ListenUDP("udp4", udpLocalAddr)
	HandleError(err)

	fmt.Println("监听", randPort, "端口")

	//用户昵称
	userName := ""
	fmt.Printf("请输入：")
	fmt.Scanf("%s", &userName)

	//_sendmes := pb.UDPPlayingGame{
	//	Userid:  proto.Int32(10001),
	//	Roomid: proto.Int32(1),
	//	Playingid: proto.Int32(2),
	//}
	_sendmes := pb.RequestCmd{
		Simple: &pb.ReqSimple{
			Tag: pb.ReqSimple_REQ_Playing_Game.Enum(),
		},
		Playinggame: &pb.UDPPlayingGame{
			Userid:  proto.Int32(10001),
			Roomid: proto.Int32(123456),
			Playingid: proto.Int32(2),
		},
	}


	senddata, err := proto.Marshal(&_sendmes)
	logs.Info("Marshal ...",&_sendmes,"Marshal:",senddata)
	//向服务器发送连接信息（昵称+本地监听端口）
	//lj :=reflectString["连接"] + userName + ":" + strconv.Itoa(randPort)
	//fmt.Println("登陆：",lj)
	//udpConn.Write([]byte(lj))


	udpConn.Write(senddata)

	//关闭端口
	defer udpConn.Close()
	defer udpListener.Close()

	//发送心跳包
	//go SendOnlineMessage(udpConn)
	//接收消息
	go HandleMessage(udpListener)

	command := ""

	for {
		//获取命令
		fmt.Printf("请输入命令：")
		fmt.Scanf("%s", &command)
		switch command {
		case "chat" :
			people := ""
			fmt.Printf("输入1：")

			fmt.Scanf("%s", &people)
			//fmt.Printf("输入您想说的话：",people)
			//向服务器发送聊天对象信息
			writerData := reflectString["聊天"]
			writerData += people
			//a := make([]int,1000)
			//for _, i2 := range a {
			//
			//	writerData += strconv.Itoa(i2)
			//}
			fmt.Println("[]byte(writerData)==>>",len([]byte(writerData)))
			udpConn.Write([]byte(writerData))
			//进入会话
			SendMessage(udpConn)
			//退出会话
			fmt.Println("退出与" + people + "的会话")


		case "get" :
			//请求在线情况信息
			udpConn.Write([]byte(reflectString["在线用户"]))
		case "1" :
			//请求在线情况信息
			udpConn.Write(senddata)
		}
	}
}

func addrss() string {

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return "192.168.3.13"
	}
	a := "192.168.3.13"
	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				a = ipnet.IP.String()
				fmt.Println("本机IP",ipnet.IP.String())
			}

		}
	}
	return a

}