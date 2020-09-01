
package tcpconn

import (
	"net"
	"os"
	_ "servicegunbattle/awesomeProject/game"
	"servicegunbattle/awesomeProject/base"
	"servicegunbattle/awesomeProject/bean"
	"servicegunbattle/awesomeProject/lib/logs"
	"servicegunbattle/awesomeProject/tcpconn/clientconn"
)
func Server() {
	//baseim.InitReadKeyMapData()
	var (
		remote = base.KeyMap.TcpAddress + ":" + base.KeyMap.Tcpport
		//data   = make([]byte, 2048)
	)
	logs.Info("TCP server start ...%v",remote)
	tcpAddr, _ := net.ResolveTCPAddr("tcp", remote)

	lis, err := net.ListenTCP("tcp",tcpAddr)

	if err != nil {
		logs.Info("Error when listen: %s, Err: %s\n", remote, err)
		os.Exit(-1)
	}else {
		logs.Info("TCP监听地址========>:%v",lis.Addr())
	}
	defer lis.Close()
	bean.Newinit()


	go clientconn.Cycle()
	//处理消息通道的逻辑

	//go clientconn.HandleMessages()
	for {
		//var res string
		conn, err := lis.Accept()
		if err != nil {
			logs.Info("Error accepting client: ", err.Error())
			os.Exit(0)
		}
		//读取客户端消息
		go clientconn.HandleConnections(conn)//data
	}
}



