package room

import (
	"net"
	"servicegunbattle/awesomeProject/pb"
)

type Player struct {
	//Id				int  		//座位id
	Userid			int  		//玩家ID
	UDPaddress		*net.UDPAddr  	//UDP 监听地址
	PlayerOpt          	*pb.PlayerOperation
}

func NewPlayer(uid int) *Player {
	s := &Player{
		//Id: id,
		Userid: uid,
		UDPaddress: &net.UDPAddr{},
		PlayerOpt: &pb.PlayerOperation{},
	}
	return s
}