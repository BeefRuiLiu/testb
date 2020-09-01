package tcpprotoapi

import (
	"bytes"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"net"
	"servicegunbattle/awesomeProject/bean"
	"servicegunbattle/awesomeProject/game"
	"servicegunbattle/awesomeProject/game/room"
	"servicegunbattle/awesomeProject/lib/logs"
	"servicegunbattle/awesomeProject/pb"
)

func RequestMatch(cmd *pb.TcpRequestMatch, cli *bean.Client,con net.Conn) {

	_room := room.NewRoom(0)
	//TODO------匹配合适room（未具体实现），没有则newroom
	game.GameMes.RoomsMap.ReadRange(func(roomid int, room *room.Room) {
		if room.Peoplecount < 2 {
			_room = room
			return
		}
	})
	logs.Info("_room.CliConnMap",_room.CliConnMap)
	l := _room.CliConnMap.Len()
	uids := make([]int32,0,l)
	clis := make([]*bean.Client,0,l)

	if _room.Roomid <= 0 {
		//创建新的room
		roomid := rand.Intn(100000) + 100000
		_room = game.GameMes.RoomsMap.CreateRoom(roomid)

		_room.RandSeed = rand.Intn(100)

	}
	s1 := room.NewPlayer(int(cmd.GetUid()))
	_room.Playermap.SetPlayer(int(cmd.GetUid()), s1)
	_room.Playermap.ForEach(func(key int, s *room.Player) {
		logs.Info("key",key,"=======>s",s)

	})

	_room.AllUIDs = append(_room.AllUIDs,int(cmd.GetUid()))

	reqUser := pb.BattleUserInfo{
		Uid: proto.Int32(cmd.GetUid()),
		BattleID:proto.Int32(int32(_room.Peoplecount+1)),
		RoleID:proto.Int32(cmd.GetRoleID()),
	}
	_room.Userinfos = append(_room.Userinfos,&reqUser)

	_room.Peoplecount += 1
	//TODO 回复TCP消息 序列化
	//isbegin := false
	//if _room.Peoplecount >= 4 {
	//	isbegin = true
	//}
	EnterBattle := &pb.TcpEnterBattle{
		RandSeed: proto.Int32(int32(_room.RandSeed)),
		BattleUserInfo:_room.Userinfos,
	}
	logs.Info("_room.Userinfos-->",_room.Userinfos)

	_room.CliConnMap.Set((cmd.GetUid()),cli)
	_room.CliConnMap.RLockRange(func(uid interface{}, clien interface{}) {
		rid := uid.(int32)
		switch clien.(type) {
		case *bean.Client:
			uids = append(uids,rid)
			clis = append(clis, clien.(*bean.Client))
		default:
			logs.Info(" ---------------> room.CliConnMap.RLockRange not found :%v", rid)
		}
	})
	if 	_room.Peoplecount < 2 {

		return
	}
	senddata, _ := proto.Marshal(EnterBattle)

	logs.Info("Marshal ...",EnterBattle,"senddata长度:",len(senddata))
	newdata := make([]byte,0,len(senddata)+3)
	var inta int8 = 50
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

	for _, _cli := range clis {
		logs.Info("clis[k]-->",_cli)
		_cli.WriteClient(newdata)
	}



	logs.Info("_room.AllUIDs ...",_room.AllUIDs)
	//写入客户端消息
	cli.WriteClient(newdata)
}
