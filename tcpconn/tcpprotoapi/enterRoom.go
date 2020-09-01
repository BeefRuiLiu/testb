package tcpprotoapi

//
//func EnterRoom(cmd *pb.RequestCmd, cli *bean.Client,con net.Conn) {
//
//	_room := &room.Room{}
//	//TODO------匹配合适room（未具体实现），没有则newroom
//	game.GameMes.RoomsMap.ReadRange(func(roomid int, room *room.Room) {
//		if room.Peoplecount < 4 {
//			_room = room
//			return
//		}
//	})
//	uids := make([]int32,0,_room.CliConnMap.Len())
//	clis := make([]bean.Client,0,_room.CliConnMap.Len())
//
//	_room.CliConnMap.RLockRange(func(uid interface{}, clien interface{}) {
//		rid := uid.(int32)
//		switch clien.(type) {
//		case *bean.Client:
//			uids = append(uids,rid)
//			clis = append(clis, clien.(bean.Client))
//		default:
//			logs.Info(" ---------------> room.CliConnMap.RLockRange not found :%v", rid)
//		}
//	})
//	if _room.Roomid <= 0 {
//		//创建新的room
//		roomid := 123456
//		_room = game.GameMes.RoomsMap.CreateRoom(roomid)
//		s1 := room.NewPlayer( cli.Uid)
//		_room.Playermap.SetPlayer(cli.Uid, s1)
//	}
//	_room.Peoplecount += 1
//	//TODO 回复TCP消息 序列化
//	isbegin := false
//	if _room.Peoplecount >= 4 {
//		isbegin = true
//	}
//		roominfo := &pb.ResponseCmd{
//			Simple: &pb.ResSimple{
//				Tag: pb.ResSimple_RES_ROOM_INFO.Enum(),
//			},
//			Roominfo:  &pb.RoomInfo{
//				Roomid: proto.Int32(int32(_room.Roomid)),
//				Count: proto.Int32(int32(_room.Peoplecount)),
//				Isbegin: proto.Bool(isbegin),
//				Userids: uids,
//			},
//		}
//		for _, _cli := range clis {
//			_cli.WriteClient(roominfo)
//		}
//
//
//	_room.CliConnMap.Set(cmd.GetEnterroom().GetUserid(),cli)
//
//	sendmes := &pb.ResponseCmd{
//		Simple: &pb.ResSimple{
//			Tag: pb.ResSimple_RES_Enter_ROOM.Enum(),
//		},
//		Resenterroom: &pb.TcpResponseEnter{
//			Roomid: proto.Int32(int32(_room.Roomid)),
//			Count: proto.Int32(int32(_room.Peoplecount)),
//			Isbegin: proto.Bool(isbegin),
//		},
//	}
//	//写入客户端消息
//	cli.WriteClient(sendmes)
//}
