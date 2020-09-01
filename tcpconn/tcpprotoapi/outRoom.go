package tcpprotoapi

//func OutRoom(cmd *pb.RequestCmd, cli *bean.Client,con net.Conn) {
//
//	//TODO------判断用户是否在房间中
//	_roomid := int(cmd.GetOutroom().GetRoomid())
//	uid := int(cmd.GetOutroom().GetUserid())
//	logs.Info("roomid:",_roomid)
//	logs.Info("uid:",uid)
//	game.GameMes.RoomsMap.ReadRange(func(roomid int, room *room.Room) {
//		if roomid == _roomid {
//			room.CliConnMap.Del(uid)
//		}
//	})
//
//	sendmes := &pb.ResponseCmd{
//		Simple: &pb.ResSimple{
//			Tag: pb.ResSimple_RES_Out_ROOM.Enum(),
//		},
//		Resoutroom: &pb.TcpResponseOut{
//			Isout: proto.Int32(1),
//		},
//	}
//	//写入客户端消息
//	cli.WriteClient(sendmes)
//
//}
